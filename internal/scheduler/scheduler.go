package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
	"weatherApi/internal/appctx"
	"weatherApi/internal/broker"
	"weatherApi/internal/common/constants"
	"weatherApi/internal/dto"
	"weatherApi/internal/logger"
	"weatherApi/internal/repository/subscription"
	serviceWeather "weatherApi/internal/service/weather"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/google/uuid"

	"github.com/go-co-op/gocron/v2"
)

const maxConcurrentJobs = 5

type SubscriptionRepositoryInterface interface {
	FindAllSubscriptionsByFrequency(ctx context.Context, frequency constants.Frequency) ([]subscription.SubscriptionModel, error)
}

type Service struct {
	log              *logger.Logger
	subscriptionRepo SubscriptionRepositoryInterface
	publisher        broker.EventPublisher
	scheduler        gocron.Scheduler
	weatherService   *serviceWeather.Service
	ctx              context.Context
}

func NewService(
	log *logger.Logger,
	subscriptionRepo SubscriptionRepositoryInterface,
	publisher broker.EventPublisher,
	weatherService *serviceWeather.Service,
	ctx context.Context,
) (*Service, error) {
	sched, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &Service{
		log:              log,
		subscriptionRepo: subscriptionRepo,
		publisher:        publisher,
		scheduler:        sched,
		weatherService:   weatherService,
		ctx:              ctx,
	}, nil
}

func (s *Service) Start() error {
	_, err := s.scheduler.NewJob(
		gocron.DurationJob(time.Hour),
		gocron.NewTask(func() {
			ctx := appctx.SetTraceID(context.Background(), uuid.NewString())
			log := s.log.FromContext(ctx)
			log.Info().Msg("Hourly job started")
			if err := s.SendNotification(ctx, constants.FrequencyHourly); err != nil {
				log.Error().Err(err).Msg("Error processing hourly notification")
			}
		}),
	)
	if err != nil {
		return err
	}

	_, err = s.scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(9, 0, 0))),
		gocron.NewTask(func() {
			ctx := appctx.SetTraceID(context.Background(), uuid.NewString())
			log := s.log.FromContext(ctx)
			log.Info().Msg("Daily job started")
			if err := s.SendNotification(ctx, constants.FrequencyDaily); err != nil {
				log.Error().Err(err).Msg("Error processing daily notification")
			}
		}),
	)
	if err != nil {
		return err
	}

	s.scheduler.Start()
	return nil
}

func (s *Service) Stop() error {
	s.log.Base().Info().Msg("Shutting down scheduler...")
	return s.scheduler.Shutdown()
}

func (s *Service) SendNotification(ctx context.Context, frequency constants.Frequency) error {
	now := time.Now()
	log := s.log.FromContext(ctx)
	if frequency == constants.FrequencyHourly && now.Hour() == 9 {
		log.Info().Msg("Hourly job skipped, cause of dayly job!")
		return nil
	}
	log.Info().Msgf("Sending notifications for %s frequency...", frequency)

	subs, err := s.subscriptionRepo.FindAllSubscriptionsByFrequency(ctx, frequency)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get subscriptions")
	}

	cityToEmails := make(map[string][]subscription.SubscriptionModel)
	for _, sub := range subs {
		city := strings.ToLower(strings.TrimSpace(sub.City))
		cityToEmails[city] = append(cityToEmails[city], sub)
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxConcurrentJobs)

	for city, subs := range cityToEmails {
		wg.Add(1)

		semaphore <- struct{}{}

		go func(ctx context.Context, city string, subs []subscription.SubscriptionModel) {
			defer wg.Done()
			defer func() { <-semaphore }()

			weather, err := s.weatherService.GetWeather(ctx, city)
			if err != nil {
				s.HandleError(fmt.Sprintf("failed to fetch weather for city=%s", city), err)
				return
			}

			users := make([]dto.UserData, len(subs))
			for i, sub := range subs {
				users[i] = dto.UserData{Email: sub.User.Email, Token: sub.ConfirmToken}
			}

			task := dto.WeatherSubData{
				Users:   users,
				Weather: *weather,
			}

			payload, marshalErr := json.Marshal(task)
			if marshalErr != nil {
				s.HandleError(fmt.Sprintf("error marshaling event for %s", city), marshalErr)
				return
			}
			traceID := appctx.GetTraceID(ctx)
			if err := s.publisher.Publish(broker.SendSubscriptionWeatherData, payload, broker.WithHeaders(amqp.Table{constants.HdrTraceID: traceID})); err != nil {
				s.HandleError(fmt.Sprintf("failed to publish notification for %s", city), err)
			}
		}(ctx, city, subs)
	}

	wg.Wait()
	return nil
}

func (s *Service) HandleError(msg string, err error) {
	s.log.Base().Error().Err(err).Msg(msg)
	if err := s.publisher.Publish(broker.SendSubscriptionWeatherData.DLQ(), []byte(msg)); err != nil {
		s.log.Base().Error().Err(err).Msg("error sending event to DLQ")
	}
}
