package main

import (
	"fmt"
	e "points-are-bad/email-server/client"

	"github.com/dangwn/points-are-bad-tooling/rmqclient"
)

func main() {
	connectionString := rmqclient.ConstructRabbitMQUrl(
		e.RABBITMQ_USER, 
		e.RABBITMQ_PASSWORD, 
		e.RABBITMQ_HOST,
		e.RABBITMQ_PORT,	
	)
	rmq, err := rmqclient.NewAsyncRabbitMQClient(connectionString, "default", e.EMAIL_VERIFICATION_QUEUE_NAME)
	if err != nil {
		fmt.Println(err)
		return
	}

	server, _ := e.NewEmailServer(
		e.AWS_ACCESS_KEY_ID,
		e.AWS_SECRET_ACCESS_KEY,
		e.AWS_REGION,
		e.AWS_SESSION_TOKEN,
		rmq,
	)
	server.Run(e.EMAIL_VERIFICATION_QUEUE_NAME)
}