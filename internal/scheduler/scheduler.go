package scheduler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
	"weatherApi/internal/broker"
	"weatherApi/internal/common/constants"
	"weatherApi/internal/common/utils"
	"weatherApi/internal/config"
	"weatherApi/internal/dto"
	"weatherApi/internal/repository/subscription"

	"github.com/go-co-op/gocron/v2"
)

const maxConcurrentJobs = 5

type Service struct {
	subscriptionRepo subscription.SubscriptionRepositoryInterface
	publisher        broker.EventPublisher
	scheduler        gocron.Scheduler
	cfg *config.Config
}

func NewService(
	subscriptionRepo subscription.SubscriptionRepositoryInterface,
	publisher broker.EventPublisher,
	cfg *config.Config,
) (*Service, error) {
	sched, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &Service{
		subscriptionRepo: subscriptionRepo,
		publisher:        publisher,
		scheduler:        sched,
		cfg: cfg,
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

			go func(city string, subs []subscription.SubscriptionModel) {
				defer wg.Done()
				defer func() { <-semaphore }()

				weather, err := s.fetchWeather(city)
				if err != nil {
					log.Printf("failed to fetch weather for city=%s: %v", city, err)
					return
				}

				users := make([]dto.UserData, len(subs))
				for i, sub := range subs {
					users[i] = dto.UserData{Email: sub.User.Email, Token: sub.ConfirmToken}
				}

				task := dto.WeatherSubData{
					Users:   users,
					Weather: weather,
				}

				payload, err := json.Marshal(task)
				if err != nil {
					log.Printf("error marshaling event for %s: %v", city, err)
					return
				}

				if err := s.publisher.Publish(broker.SendSubscriptionWeatherData, payload); err != nil {
					log.Printf("failed to publish notification for %s: %v", city, err)
				}
			}(city, subs)
		}

		wg.Wait()
		return nil
	}

	func (s *Service) fetchWeather(city string) (dto.WeatherResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var weatherResponse dto.WeatherResponse

		payload, err := json.Marshal(map[string]string{
			"city": city,
		})
		if err != nil {
			return weatherResponse, fmt.Errorf("failed to marshal city: %w", err)
		}

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			fmt.Sprintf("%s/api/v1/weather", s.cfg.AppURL),
			bytes.NewBuffer(payload),
		)
		if err != nil {
			return weatherResponse, fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		response, err := http.DefaultClient.Do(req)
		if err != nil {
			return weatherResponse, fmt.Errorf("request failed: %w", err)
		}

		defer func() {
			if err := response.Body.Close(); err != nil {
				log.Printf("failed to close response body: %v", err)
			}
		}()

		if response.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(response.Body)
			return weatherResponse, fmt.Errorf("unexpected status %d: %s", response.StatusCode, string(body))
		}

		if err := json.NewDecoder(response.Body).Decode(&weatherResponse); err != nil {
			return weatherResponse, fmt.Errorf("failed to decode response: %w", err)
		}

		return weatherResponse, nil
	}

	func (s *Service) fetchWeatherWithRetry(city string) (dto.WeatherResponse, error) {
		return utils.Retry[dto.WeatherResponse](3, 2*time.Second, func() (dto.WeatherResponse, error) {
			return s.fetchWeather(city)
		})
	}