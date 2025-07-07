package subscription

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"
	"weatherApi/internal/broker"
	"weatherApi/internal/repository/base"

	"weatherApi/internal/common/constants"
	commonErrors "weatherApi/internal/common/errors"
	"weatherApi/internal/dto"
	"weatherApi/internal/repository/subscription"
	"weatherApi/internal/repository/user"
	serviceErrors "weatherApi/internal/service/subscription/errors"

	"github.com/google/uuid"
)

type SubscriptionService struct {
	SubscriptionRepo subscription.SubscriptionRepositoryInterface
	UserRepo         user.UserRepositoryInterface
	publisher        broker.EventPublisher
	tokenLifeMinutes int
}

func NewSubscriptionService(
	subscriptionRepo subscription.SubscriptionRepositoryInterface,
	userRepo user.UserRepositoryInterface,
	publisher broker.EventPublisher,
	tokenLifeMinutes int,
) *SubscriptionService {
	return &SubscriptionService{
		SubscriptionRepo: subscriptionRepo,
		UserRepo:         userRepo,
		publisher:        publisher,
		tokenLifeMinutes: tokenLifeMinutes,
	}
}

func (s *SubscriptionService) Subscribe(ctx context.Context, subscribeRequest *dto.SubscribeRequest) *commonErrors.AppError {
	token, err := s.generateConfirmationToken()
	if err != nil {
		return serviceErrors.ErrInternalServerError
	}

	userModel := &user.UserModel{
		Email: subscribeRequest.Email,
	}

	user, err := s.UserRepo.FindOneOrCreate(ctx, map[string]any{
		"email": subscribeRequest.Email,
	}, userModel)
	if err != nil {
		return serviceErrors.ErrInternalServerError
	}

	expiry := time.Now().Add(time.Duration(s.tokenLifeMinutes) * time.Minute)

	existing, err := s.SubscriptionRepo.FindOneOrNone(ctx, "user_id = ?", user.ID)
	if err != nil {
		if errors.Is(err, base.ErrNotFound) {
			existing = &subscription.SubscriptionModel{
				City:         subscribeRequest.City,
				Frequency:    constants.Frequency(subscribeRequest.Frequency),
				UserID:       user.ID,
				IsConfirmed:  false,
				ConfirmToken: token,
				TokenExpires: expiry,
			}

			if err := s.SubscriptionRepo.CreateOne(ctx, existing); err != nil {
				return serviceErrors.ErrInternalServerError
			}
		} else {
			return serviceErrors.ErrInternalServerError
		}
	}

	if existing.IsConfirmed {
		return serviceErrors.ErrAlreadySubscribed
	}

	existing.ConfirmToken = token
	existing.TokenExpires = expiry
	existing.Frequency = constants.Frequency(subscribeRequest.Frequency)

	if err := s.SubscriptionRepo.Update(ctx, existing); err != nil {
		return serviceErrors.ErrInternalServerError
	}

	task := dto.ConfirmationEmailTask{
		Email: subscribeRequest.Email,
		Token: token,
		City:  subscribeRequest.City,
	}
	payload, err := json.Marshal(task)
	if err != nil {
		log.Println("Error marshaling confirmation event")
		return serviceErrors.ErrInternalServerError
	}
	if err := s.publisher.Publish(broker.SubscriptionConfirmationTasks, payload); err != nil {
		log.Println("Error publishing confirmation event")
		return serviceErrors.ErrInternalServerError
	}
	log.Printf("New task published! %s", task.Email)
	return nil
}

func (s *SubscriptionService) ConfirmSubscription(ctx context.Context, token string) *commonErrors.AppError {
	subscription, err := s.SubscriptionRepo.FindOneOrNone(
		ctx,
		"confirm_token = ? AND deleted_at IS NULL",
		token,
	)
	if err != nil {
		if errors.Is(err, base.ErrNotFound) {
			return serviceErrors.ErrTokenNotFound
		}
		return serviceErrors.ErrInternalServerError
	}

	if subscription.IsConfirmed {
		return serviceErrors.ErrAlreadySubscribed
	}

	if time.Now().After(subscription.TokenExpires) {
		return serviceErrors.ErrInvalidToken
	}
	now := time.Now()
	subscription.IsConfirmed = true
	subscription.ConfirmedAt = &now

	if err := s.SubscriptionRepo.Update(ctx, subscription); err != nil {
		return serviceErrors.ErrInternalServerError
	}

	return nil
}

func (s *SubscriptionService) Unsubscribe(ctx context.Context, token string) *commonErrors.AppError {
	subscription, err := s.SubscriptionRepo.FindOneOrNone(
		ctx,
		"confirm_token = ? AND is_confirmed = ?",
		token,
		true,
	)
	if err != nil {
		if errors.Is(err, base.ErrNotFound) {
			return serviceErrors.ErrTokenNotFound
		}
		return serviceErrors.ErrInternalServerError
	}

	if err := s.SubscriptionRepo.Delete(ctx, subscription); err != nil {
		return serviceErrors.ErrInternalServerError
	}

	return nil
}

func (s *SubscriptionService) generateConfirmationToken() (string, error) {
	token, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return token.String(), nil
}
