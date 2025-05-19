package provider

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type SMTPClientInterface interface {
	SendConfirmationToken(to, token, city string) error
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
