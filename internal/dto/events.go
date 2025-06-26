package dto

type ConfirmationEmailTask struct {
	Email string `json:"email"`
	Token string `json:"token"`
	City  string `json:"city"`
}
