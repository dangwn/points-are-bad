package rmqclient

import (
	"fmt"
	"log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AsyncRabbitMQClient struct {
	Connection *amqp.Connection
	Channel *amqp.Channel
	Exchange string
	Queues map[string]amqp.Queue
	Messages <-chan amqp.Delivery
	Done chan error
}

func NewAsyncRabbitMQClient(
	connectionUrl string, 
	exchangeName string, 
	queueNames []string,
) (*AsyncRabbitMQClient) {
	conn, err := amqp.Dial(connectionUrl)
	if err != nil {
		log.Fatal(err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		log.Fatal(err)
	}

	queues := make(map[string]amqp.Queue)
	for _, name := range queueNames {
		q, err := channel.QueueDeclare(
			name, // name
			true,   // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		if err != nil {
			conn.Close()
			channel.Close()
			log.Fatal(err)
		}
		queues[name] = q
	}

	if exchangeName != "" {	
		// If the exchange name is undefined, use the default exchange
		err = channel.ExchangeDeclare(
			exchangeName, // name
			"direct", // type
			true,     // durable
			false,    // auto-deleted
			false,    // internal
			false,    // no-wait
			nil,      // arguments
		)

		if err != nil {
			conn.Close()
			channel.Close()
			log.Fatal(err)
		}

		// Bind queues to new exchange
		for _, name := range queueNames {
			err = channel.QueueBind(
				queues[name].Name,
				name,
				exchangeName,
				false, // no wait
				nil, // arguments
			)

			if err != nil {
				conn.Close()
				channel.Close()
				log.Fatal(err)
			}
		}
	}

	return &AsyncRabbitMQClient{
		Connection: conn,
		Channel:	channel,
		Exchange:	exchangeName,
		Queues:     queues,
	}
}

func (c *AsyncRabbitMQClient) Shutdown() {
	c.Connection.Close()
	c.Channel.Close()
}

func ConstructRabbitMQUrl(
	username string, 
	password string, 
	hostname string,
	port string,
) string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		username,
		password,
		hostname,
		port,
	)
}