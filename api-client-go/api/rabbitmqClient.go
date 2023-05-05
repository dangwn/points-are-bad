package api

import (
	rmq "points-are-bad/api-client/rabbitmq-client"
)

var rabbit rmq.AsyncRabbitMQClient = *rmq.NewAsyncRabbitMQClient(
	rmq.ConstructRabbitMQUrl(
		RABBITMQ_USER,
		RABBITMQ_PASSWORD,
		RABBITMQ_HOST,
		RABBITMQ_PORT,
	),
	"emailExchange",
	[]string{EMAIL_VERIFICATION_QUEUE},
)