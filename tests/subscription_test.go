package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"weatherApi/internal/broker"
	"weatherApi/internal/logger"
	"weatherApi/internal/repository/base"

	"weatherApi/internal/repository/subscription"
	"weatherApi/internal/repository/user"
	"weatherApi/internal/server/routes"
	subscriptionService "weatherApi/internal/service/subscription"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter(handler *routes.SubscriptionHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/subscribe", handler.Subscribe)
	r.GET("/confirm/:token", handler.ConfirmSubscription)
	r.GET("/unsubscribe/:token", handler.Unsubscribe)
	return r
}

func TestSubscribeSuccess(t *testing.T) {
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
	log := logger.NewNoOpLogger()

	service := subscriptionService.NewSubscriptionService(log, subRepo, userRepo, publisher, 60)
	handler := routes.NewSubscriptionHandler(log, service)
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
}

func TestSubscribeInvalidInput(t *testing.T) {
	log := logger.NewNoOpLogger()
	service := subscriptionService.NewSubscriptionService(log, nil, nil, nil, 60)
	handler := routes.NewSubscriptionHandler(log, service)
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
}

func TestConfirmSubscriptionSuccess(t *testing.T) {
	sub := &subscription.SubscriptionModel{
		IsConfirmed:  false,
		TokenExpires: time.Now().Add(1 * time.Hour),
	}
	subRepo := &subscription.MockSubscriptionRepository{
		FindOneOrNoneFn: func(_ any, _ ...any) (*subscription.SubscriptionModel, error) {
			return sub, nil
		},
		UpdateFn: func(_ *subscription.SubscriptionModel) error {
			return nil
		},
	}
	log := logger.NewNoOpLogger()
	service := subscriptionService.NewSubscriptionService(log, subRepo, nil, nil, 60)
	handler := routes.NewSubscriptionHandler(log, service)
	router := setupTestRouter(handler)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/confirm/some-valid-token", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Subscription confirmed")
}

func TestConfirmSubscriptionInvalidToken(t *testing.T) {
	subRepo := &subscription.MockSubscriptionRepository{
		FindOneOrNoneFn: func(_ any, _ ...any) (*subscription.SubscriptionModel, error) {
			return nil, base.ErrNotFound
		},
	}
	log := logger.NewNoOpLogger()
	service := subscriptionService.NewSubscriptionService(log, subRepo, nil, nil, 60)
	handler := routes.NewSubscriptionHandler(log, service)
	router := setupTestRouter(handler)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/confirm/some-valid-token", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Token not found")
}

func TestTokenExpired(t *testing.T) {
	sub := &subscription.SubscriptionModel{
		IsConfirmed:  false,
		TokenExpires: time.Now().Add(-1 * time.Hour),
	}
	subRepo := &subscription.MockSubscriptionRepository{
		FindOneOrNoneFn: func(_ any, _ ...any) (*subscription.SubscriptionModel, error) {
			return sub, nil
		},
		UpdateFn: func(_ *subscription.SubscriptionModel) error {
			return nil
		},
	}
	log := logger.NewNoOpLogger()
	service := subscriptionService.NewSubscriptionService(log, subRepo, nil, nil, 60)
	handler := routes.NewSubscriptionHandler(log, service)
	router := setupTestRouter(handler)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/confirm/some-valid-token", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token")
}

func TestUnsubscribeSuccess(t *testing.T) {
	sub := &subscription.SubscriptionModel{
		IsConfirmed: true,
	}
	subRepo := &subscription.MockSubscriptionRepository{
		FindOneOrNoneFn: func(_ any, _ ...any) (*subscription.SubscriptionModel, error) {
			return sub, nil
		},
		DeleteFn: func(_ *subscription.SubscriptionModel) error {
			return nil
		},
	}
	log := logger.NewNoOpLogger()
	service := subscriptionService.NewSubscriptionService(log, subRepo, nil, nil, 60)
	handler := routes.NewSubscriptionHandler(log, service)
	router := setupTestRouter(handler)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/unsubscribe/some-token", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Unsubscribed successfully")
}

func TestUnsubscribeSubscriptionTokenNotFound(t *testing.T) {
	subRepo := &subscription.MockSubscriptionRepository{
		FindOneOrNoneFn: func(_ any, _ ...any) (*subscription.SubscriptionModel, error) {
			return nil, base.ErrNotFound
		},
	}
	log := logger.NewNoOpLogger()
	service := subscriptionService.NewSubscriptionService(log, subRepo, nil, nil, 60)
	handler := routes.NewSubscriptionHandler(log, service)
	router := setupTestRouter(handler)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/unsubscribe/some-token", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Token not found")
}
