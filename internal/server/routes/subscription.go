package routes

import (
	"net/http"
	"weatherApi/internal/logger"

	"weatherApi/internal/dto"
	"weatherApi/internal/service/subscription"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	log     *logger.Logger
	service *subscription.SubscriptionService
}

func NewSubscriptionHandler(
	log *logger.Logger,
	subscriptionService *subscription.SubscriptionService,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		log:     log,
		service: subscriptionService,
	}
}

func (h *SubscriptionHandler) Subscribe(c *gin.Context) {
	log := h.log.FromContext(c.Request.Context())
	var req dto.SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err := h.service.Subscribe(c.Request.Context(), &req); err != nil {
		log.Error().Err(err).Msgf("Failed to handle subscribe request for %s: %s", req.Email, req.City)
		c.AbortWithStatusJSON(err.Code, gin.H{"error": err.Message})
		return
	}

	c.JSON(http.StatusOK, "Subscription successful. Confirmation email sent.")
}

func (h *SubscriptionHandler) ConfirmSubscription(c *gin.Context) {
	log := h.log.FromContext(c.Request.Context())
	token := c.Param("token")
	log.Info().Msgf("Handling confirm subscription for token %s", token)

	if err := h.service.ConfirmSubscription(c.Request.Context(), token); err != nil {
		log.Error().Err(err).Msgf("Failed to handle confirm subscription for token %s", token)
		c.AbortWithStatusJSON(err.Code, gin.H{"error": err.Message})
		return
	}
	c.JSON(http.StatusOK, "Subscription confirmed successfully")
}

func (h *SubscriptionHandler) Unsubscribe(c *gin.Context) {
	log := h.log.FromContext(c.Request.Context())
	token := c.Param("token")
	log.Info().Msgf("Handling unsibscribe from subscription for token %s", token)
	if err := h.service.Unsubscribe(c.Request.Context(), token); err != nil {
		log.Error().Err(err).Msgf("Failed to handle unsubscribe for token %s", token)
		c.AbortWithStatusJSON(err.Code, gin.H{"error": err.Message})
		return
	}
	c.JSON(http.StatusOK, "Unsubscribed successfully")
}
