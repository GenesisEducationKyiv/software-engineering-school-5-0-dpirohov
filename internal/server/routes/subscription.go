package routes

import (
	"net/http"

	"weatherApi/internal/service/subscription"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	service *subscription.SubscriptionService
}

type SubscribeRequest struct {
	Email     string `json:"email"     binding:"required"`
	City      string `json:"city"      binding:"required"`
	Frequency string `json:"frequency" binding:"required,oneof=hourly daily"`
}

func NewSubscriptionHandler(
	subscriptionService *subscription.SubscriptionService,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: subscriptionService,
	}
}

func (h *SubscriptionHandler) Subscribe(c *gin.Context) {
	var req SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.service.Subscribe(req.Email, req.City, req.Frequency); err != nil {
		c.AbortWithStatusJSON(err.Code, gin.H{"error": err.Message})
		return
	}

	c.JSON(http.StatusOK, "Subscription successful. Confirmation email sent.")
}

func (h *SubscriptionHandler) ConfirmSubscription(c *gin.Context) {
	token := c.Param("token")
	if err := h.service.ConfirmSubscription(token); err != nil {
		c.AbortWithStatusJSON(err.Code, gin.H{"error": err.Message})
		return
	}
	c.JSON(http.StatusOK, "Subscription confirmed successfully")
}

func (h *SubscriptionHandler) Unsubscribe(c *gin.Context) {
	token := c.Param("token")
	if err := h.service.Unsubscribe(token); err != nil {
		c.AbortWithStatusJSON(err.Code, gin.H{"error": err.Message})
		return
	}
	c.JSON(http.StatusOK, "Unsubscribed successfully")
}
