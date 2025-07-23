package dto

type ConfirmationEmailTask struct {
	Email string `json:"email"`
	Token string `json:"token"`
	City  string `json:"city"`
}


type UserData struct {
	Email string `json:"Email"`
	Token string `json:"Token"`
}


type WeatherSubData struct {
	Users []UserData `json:"users"`
	Weather WeatherResponse `json:"weather"`
}