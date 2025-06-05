package subscription

import (
	"time"

	"weatherApi/internal/common/constants"
	commonErrors "weatherApi/internal/common/errors"
	"weatherApi/internal/provider"
	"weatherApi/internal/repository/subscription"
	"weatherApi/internal/repository/user"
	serviceErrors "weatherApi/internal/service/subscription/errors"

	"github.com/google/uuid"
)

type SubscriptionService struct {
	SubscriptionRepo subscription.SubscriptionRepositoryInterface
	UserRepo         user.UserRepositoryInterface
	smtpClient       provider.SMTPClientInterface
	tokenLifeMinutes int
}

func NewSubscriptionService(
	subscriptionRepo subscription.SubscriptionRepositoryInterface,
	userRepo user.UserRepositoryInterface,
	smtpClient provider.SMTPClientInterface,
	tokenLifeMinutes int,
) *SubscriptionService {
	return &SubscriptionService{
		SubscriptionRepo: subscriptionRepo,
		UserRepo:         userRepo,
		smtpClient:       smtpClient,
		tokenLifeMinutes: tokenLifeMinutes,
	}
}

func (s *SubscriptionService) Subscribe(email, city, frequency string) *commonErrors.AppError {
	if email == "" || city == "" || frequency == "" {
		return serviceErrors.ErrInternalServerError
	}

	token, err := s.generateConfirmationToken()
	if err != nil {
		return serviceErrors.ErrInternalServerError
	}

	userModel := &user.UserModel{
		Email: email,
	}

	user, err := s.UserRepo.FindOneOrCreate(map[string]any{
		"email": email,
	}, userModel)
	if err != nil {
		return serviceErrors.ErrInternalServerError
	}

	existing, err := s.SubscriptionRepo.FindOneOrNone("user_id = ?", user.ID)
	if err != nil {
		return serviceErrors.ErrInternalServerError
	}

	expiry := time.Now().Add(time.Duration(s.tokenLifeMinutes) * time.Minute)

	if existing != nil {
		if existing.IsConfirmed {
			return serviceErrors.ErrAlreadySubscribed
		}

		existing.ConfirmToken = token
		existing.TokenExpires = expiry
		existing.Frequency = constants.Frequency(frequency)

		if err := s.SubscriptionRepo.Update(existing); err != nil {
			return serviceErrors.ErrInternalServerError
		}
	} else {
		newSub := &subscription.SubscriptionModel{
			City:         city,
			Frequency:    constants.Frequency(frequency),
			UserID:       user.ID,
			IsConfirmed:  false,
			ConfirmToken: token,
			TokenExpires: expiry,
		}

		if err := s.SubscriptionRepo.CreateOne(newSub); err != nil {
			return serviceErrors.ErrInternalServerError
		}
	}

	if err := s.smtpClient.SendConfirmationToken(email, token, city); err != nil {
		return serviceErrors.ErrInternalServerError
	}

	return nil
}

func (s *SubscriptionService) ConfirmSubscription(token string) *commonErrors.AppError {
	subscription, err := s.SubscriptionRepo.FindOneOrNone(
		"confirm_token = ? AND deleted_at IS NULL",
		token,
	)
	if err != nil {
		return serviceErrors.ErrInternalServerError
	}
	if subscription == nil {
		return serviceErrors.ErrTokenNotFound
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

	if err := s.SubscriptionRepo.Update(subscription); err != nil {
		return serviceErrors.ErrInternalServerError
	}

	return nil
}

func (s *SubscriptionService) Unsubscribe(token string) *commonErrors.AppError {
	subscription, err := s.SubscriptionRepo.FindOneOrNone(
		"confirm_token = ? AND is_confirmed = ?",
		token,
		true,
	)
	if err != nil {
		return serviceErrors.ErrInternalServerError
	}
	if subscription == nil {
		return serviceErrors.ErrTokenNotFound
	}

	if err := s.SubscriptionRepo.Delete(subscription); err != nil {
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
