package broker

import "strings"

type Topic string

const (
	SubscriptionConfirmationTasks Topic = "task.send_confirmation_token"
	SendSubscriptionWeatherData   Topic = "task.send_sub_data"
)

func (t Topic) DLQ() Topic {
	return Topic(strings.Replace(string(t), "task", "dlq", 1))
}
