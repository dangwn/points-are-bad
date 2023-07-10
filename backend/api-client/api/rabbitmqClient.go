package api

import (
	rmq "github.com/dangwn/points-are-bad-tooling/rmqclient"
)

var rabbit rmq.AsyncRabbitMQClient = func() rmq.AsyncRabbitMQClient {
	url := rmq.ConstructRabbitMQUrl(RABBITMQ_USER, RABBITMQ_PASSWORD, RABBITMQ_HOST,RABBITMQ_PORT)
	client, err := rmq.NewAsyncRabbitMQClient(url, "emailExchange", EMAIL_VERIFICATION_QUEUE)
	if err != nil {
		panic("could not start rabbitmq client due to following error: " + err.Error())
	}
	return *client
}()