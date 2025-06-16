package broker

type Topic string

const (
	SubscriptionConfirmationTasks Topic = "task.send_confirmation_token"
	DeadLetterQueue               Topic = "dlq.all"
)

var AllTopics = []Topic{
	SubscriptionConfirmationTasks,
	DeadLetterQueue,
}
