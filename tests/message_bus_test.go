package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"weatherApi/internal/broker"
	"weatherApi/internal/common/constants"
	"weatherApi/internal/dto"
	"weatherApi/internal/repository/base"

	"weatherApi/internal/repository/subscription"
	"weatherApi/internal/repository/user"
	"weatherApi/internal/server/routes"
	subscriptionService "weatherApi/internal/service/subscription"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPublisherCalled(t *testing.T) {
	userRepo := &user.MockUserRepository{
		FindOneOrCreateFn: func(_ map[string]any, e *user.UserModel) (*user.UserModel, error) {
			e.ID = 1
			return e, nil
		},
	}

	subRepo := &subscription.MockSubscriptionRepository{
		FindOneOrNoneFn: func(_ any, _ ...any) (*subscription.SubscriptionModel, error) {
			return nil, base.ErrNotFound
		},
		CreateOneFn: func(_ *subscription.SubscriptionModel) error {
			return nil
		},
		UpdateFn: func(entity *subscription.SubscriptionModel) error {
			return nil
		},
	}

	publisher := broker.NewMockRabbitMQPublisher()

	service := subscriptionService.NewSubscriptionService(subRepo, userRepo, publisher, 60)
	handler := routes.NewSubscriptionHandler(service)
	router := setupTestRouter(handler)

	body, _ := json.Marshal(gin.H{
		"email":     "test@example.com",
		"city":      "Kyiv",
		"frequency": "daily",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/subscribe", bytes.NewBuffer(body))

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Subscription successful")
	assert.Len(t, publisher.Calls, 1)
	assert.Equal(t, broker.SubscriptionConfirmationTasks, publisher.Calls[0].Topic)
	var actual dto.ConfirmationEmailTask
	if err := json.Unmarshal(publisher.Calls[0].Payload, &actual); err != nil {
		t.Fail()
	}
	assert.Equal(t, "test@example.com", actual.Email)
	assert.Equal(t, "Kyiv", actual.City)
}

func TestPublisherNotCalledOnInvalidInput(t *testing.T) {
	userRepo := &user.MockUserRepository{
		FindOneOrCreateFn: func(_ map[string]any, e *user.UserModel) (*user.UserModel, error) {
			e.ID = 1
			return e, nil
		},
	}

	subRepo := &subscription.MockSubscriptionRepository{
		FindOneOrNoneFn: func(_ any, _ ...any) (*subscription.SubscriptionModel, error) {
			return nil, base.ErrNotFound
		},
		CreateOneFn: func(_ *subscription.SubscriptionModel) error {
			return nil
		},
		UpdateFn: func(entity *subscription.SubscriptionModel) error {
			return nil
		},
	}

	publisher := broker.NewMockRabbitMQPublisher()

	service := subscriptionService.NewSubscriptionService(subRepo, userRepo, publisher, 60)
	handler := routes.NewSubscriptionHandler(service)
	router := setupTestRouter(handler)

	body, _ := json.Marshal(gin.H{
		"email": "test@example.com", // missing "city" and "frequency"
	})
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/subscribe", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid input")
	assert.Len(t, publisher.Calls, 0)
}

func TestPublisherNotCalledOnInternalError(t *testing.T) {
	userRepo := &user.MockUserRepository{
		FindOneOrCreateFn: func(_ map[string]any, e *user.UserModel) (*user.UserModel, error) {
			e.ID = 1
			return e, nil
		},
	}

	subRepo := &subscription.MockSubscriptionRepository{
		FindOneOrNoneFn: func(_ any, _ ...any) (*subscription.SubscriptionModel, error) {
			return nil, fmt.Errorf("mock error")
		},
		CreateOneFn: func(_ *subscription.SubscriptionModel) error {
			return nil
		},
		UpdateFn: func(entity *subscription.SubscriptionModel) error {
			return nil
		},
	}

	publisher := broker.NewMockRabbitMQPublisher()

	service := subscriptionService.NewSubscriptionService(subRepo, userRepo, publisher, 60)
	handler := routes.NewSubscriptionHandler(service)
	router := setupTestRouter(handler)

	body, _ := json.Marshal(gin.H{
		"email":     "test@example.com",
		"city":      "Kyiv",
		"frequency": "daily",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/subscribe", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Internal server error")
	assert.Len(t, publisher.Calls, 0)
}

func TestPublisherNotCalled_SubscriptionAlreadyConfirmed(t *testing.T) {
	userRepo := &user.MockUserRepository{
		FindOneOrCreateFn: func(_ map[string]any, e *user.UserModel) (*user.UserModel, error) {
			e.ID = 1
			return e, nil
		},
	}

	subRepo := &subscription.MockSubscriptionRepository{
		FindOneOrNoneFn: func(_ any, _ ...any) (*subscription.SubscriptionModel, error) {
			return &subscription.SubscriptionModel{
				City:         "Test",
				Frequency:    constants.FrequencyHourly,
				UserID:       1,
				User:         user.UserModel{},
				IsConfirmed:  true,
				ConfirmToken: "test",
			}, nil
		},
		CreateOneFn: func(_ *subscription.SubscriptionModel) error {
			return nil
		},
		UpdateFn: func(entity *subscription.SubscriptionModel) error {
			return nil
		},
	}

	publisher := broker.NewMockRabbitMQPublisher()

	service := subscriptionService.NewSubscriptionService(subRepo, userRepo, publisher, 60)
	handler := routes.NewSubscriptionHandler(service)
	router := setupTestRouter(handler)

	body, _ := json.Marshal(gin.H{
		"email":     "test@example.com",
		"city":      "Kyiv",
		"frequency": "daily",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/subscribe", bytes.NewBuffer(body))

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "mail already subscribed")
	assert.Len(t, publisher.Calls, 0)
}
