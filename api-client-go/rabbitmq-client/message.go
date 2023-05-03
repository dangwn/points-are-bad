package rmqclient

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (c *AsyncRabbitMQClient) SendMessage(
	body string,
	queueName string,
	timeoutSeconds int,
) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(timeoutSeconds) * time.Second,
	)
	defer cancel()
	
	err := c.Channel.PublishWithContext(
		ctx,
		"", // exchange
		c.Queues[queueName].Name, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:		  []byte(body),
		},
	)
	return err
}

func (c *AsyncRabbitMQClient) ConsumeMessages(
	queueName string,
	outChannel chan string,
	stopConsumingChannel chan bool,
) {
	messages, err := c.Channel.Consume(
		c.Queues[queueName].Name,
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return
	}
	for {
		select {
		case msg := <- messages:
			msgBody := string(msg.Body)
			outChannel <- msgBody
		case <- stopConsumingChannel:
			return
		}		
	}
}