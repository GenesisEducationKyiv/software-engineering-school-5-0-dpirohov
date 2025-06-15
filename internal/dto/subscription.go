package dto

type SubscribeRequest struct {
	Email     string `json:"email"     binding:"required,email"`
	City      string `json:"city"      binding:"required"`
	Frequency string `json:"frequency" binding:"required,oneof=hourly daily"`
}
