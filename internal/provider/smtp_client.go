package provider

import (
	"fmt"
	"weatherApi/internal/dto"

	"gopkg.in/gomail.v2"
)

type SMTPClientInterface interface {
	SendConfirmationToken(to, token, city string) error
	SendSubscriptionWeatherData(data *dto.WeatherResponse, user *dto.UserData) error
}

type SMTPClient struct {
	host      string
	port      int
	login     string
	password  string
	serverUrl string
}

func NewSMTPClient(host string, port int, login, password, serverUrl string) SMTPClientInterface {
	return &SMTPClient{
		host:      host,
		port:      port,
		login:     login,
		password:  password,
		serverUrl: serverUrl,
	}
}

func (c *SMTPClient) SendConfirmationToken(to, token, city string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", c.login)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Weather subscription confimation")
	confirmationURL := fmt.Sprintf("%s/confirm/%s", c.serverUrl, token)
	htmlBody := fmt.Sprintf(`
		<html>
			<body style="font-family: Arial, sans-serif; color: #333;">
				<h2>Hello!</h2>
				<p>You requested to subscribe to weather updates for <strong>%s</strong>.</p>
				<p>Please confirm your subscription by clicking the button below:</p>
				<a href="%s" 
				   style="display:inline-block; padding:10px 20px; background-color:#28a745; color:white; text-decoration:none; border-radius:5px;">
					Confirm Subscription
				</a>
				<p>If you did not request this, you can ignore this email.</p>
				<br/>
				<small>Weather Service Team</small>
			</body>
		</html>`, city, confirmationURL)

	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(c.host, c.port, c.login, c.password)

	return d.DialAndSend(m)
}

func (c *SMTPClient) SendSubscriptionWeatherData(data *dto.WeatherResponse, user *dto.UserData) error {
	m := gomail.NewMessage()
	m.SetHeader("From", c.login)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Weather subscription confimation")
	unsubscribeUrl := fmt.Sprintf("%s/unsubscribe/%s", c.serverUrl, user.Token)
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Weather Forecast</title>
  <style>
    body {
      font-family: Arial, sans-serif;
      background-color: #f4f6f8;
      padding: 20px;
      color: #333;
    }
    .container {
      background-color: #ffffff;
      border-radius: 8px;
      padding: 24px;
      max-width: 500px;
      margin: auto;
      box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }
    .heading {
      font-size: 22px;
      font-weight: bold;
      margin-bottom: 16px;
      text-align: center;
    }
    .info {
      font-size: 16px;
      margin-bottom: 10px;
    }
    .footer {
      margin-top: 20px;
      font-size: 12px;
      color: #888;
      text-align: center;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="heading">🌤️ Weather Update</div>
    <div class="info">🌡️ <strong>Temperature:</strong> %.1f°C</div>
    <div class="info">💧 <strong>Humidity:</strong> %d%%</div>
    <div class="info">📖 <strong>Conditions:</strong> %d</div>
    <div class="footer">You are receiving this weather update because you subscribed to weather notifications.</div>
    <div class="unsubscribe">
      👉 <a href="%s">Unsubscribe from future updates</a>
    </div>  </div>
</body>
</html>
`, data.Temperature, data.Humidity, data.Humidity, unsubscribeUrl)

	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(c.host, c.port, c.login, c.password)

	return d.DialAndSend(m)
}
