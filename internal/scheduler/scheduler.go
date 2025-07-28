package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
	"weatherApi/internal/broker"
	"weatherApi/internal/common/constants"
	"weatherApi/internal/dto"
	"weatherApi/internal/repository/subscription"
	serviceWeather "weatherApi/internal/service/weather"

	"github.com/go-co-op/gocron/v2"
)

const maxConcurrentJobs = 5

type Service struct {
	subscriptionRepo subscription.SubscriptionRepositoryInterface
	publisher        broker.EventPublisher
	scheduler        gocron.Scheduler
	weatherService   *serviceWeather.Service
	ctx              context.Context
}

func NewService(
	subscriptionRepo subscription.SubscriptionRepositoryInterface,
	publisher broker.EventPublisher,
	weatherService *serviceWeather.Service,
	ctx context.Context,
) (*Service, error) {
	sched, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &Service{
		subscriptionRepo: subscriptionRepo,
		publisher:        publisher,
		scheduler:        sched,
		weatherService:   weatherService,
		ctx:              ctx,
	}, nil
}

func (s *Service) Start(ctx context.Context) error {
	_, err := s.scheduler.NewJob(
		gocron.DurationJob(time.Hour),
		gocron.NewTask(func() {
			log.Println("Hourly job started")
			if err := s.SendNotification(ctx, constants.FrequencyHourly); err != nil {
				log.Printf("error in hourly notification: %v", err)
			}
		}),
	)
	if err != nil {
		return err
	}

	_, err = s.scheduler.NewJob(
		gocron.CronJob("0 9 * * *", false),
		gocron.NewTask(func() {
			log.Println("Daily job started")
			if err := s.SendNotification(ctx, constants.FrequencyDaily); err != nil {
				log.Printf("error in daily notification: %v", err)
			}
		}),
	)
	if err != nil {
		return err
	}

	s.scheduler.Start()
	return nil
}

func (s *Service) Stop(ctx context.Context) error {
	log.Println("Shutting down scheduler...")
	return s.scheduler.Shutdown()
}

func (s *Service) SendNotification(ctx context.Context, frequency constants.Frequency) error {
	now := time.Now()

	if frequency == constants.FrequencyHourly && now.Hour() == 9 {
		log.Println("Hourly job skipped, cause of dayly job!")
		return nil
	}

	log.Printf("Sending notifications for %s frequency...\n", frequency)

	subs, err := s.subscriptionRepo.FindAllSubscriptionsByFrequency(ctx, frequency)
	if err != nil {
		return fmt.Errorf("failed to get subscriptions: %w", err)
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
				s.HandleError(fmt.Sprintf("failed to fetch weather for city=%s: %v", city, err))
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
				s.HandleError(fmt.Sprintf("error marshaling event for %s: %v", city, marshalErr))
				return
			}

			if err := s.publisher.Publish(broker.SendSubscriptionWeatherData, payload); err != nil {
				s.HandleError(fmt.Sprintf("failed to publish notification for %s: %v", city, err))
			}
		}(ctx, city, subs)
	}

	wg.Wait()
	return nil
}

func (s *Service) HandleError(msg string) {
	log.Println(msg)
	if err := s.publisher.Publish(broker.SendSubscriptionWeatherData.DLQ(), []byte(msg)); err != nil {
		log.Printf("error sending event to DLQ %v", err)
	}
}
