package rmqclient

import (
	"context"
	"time"

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

/*
 * Creates a new rabbitmq client, and cleans up connections and channels if errors are thrown
 */
func NewAsyncRabbitMQClient(connectionUrl, exchangeName string, queueNames []string) (a *AsyncRabbitMQClient, e error) {
	var conn *amqp.Connection
	var channel *amqp.Channel
	var err error

	// Cleanup of connections and channels
	defer func() {
		if r := recover(); r != nil {
			if conn != nil {
				conn.Close()
			}
			if channel != nil {
				channel.Close()
			}
			a = nil
			e = r.(error)
		}
	}()	

	if conn, err = amqp.Dial(connectionUrl); err != nil {
		panic(err)
	}

	if channel, err = conn.Channel(); err != nil {
		panic(err)
	}

	queues := make(map[string]amqp.Queue)
	for _, name := range queueNames {
		if q, err := channel.QueueDeclare(name, true, false, false, false, nil); err != nil {
			panic(err)
		} else {
			queues[name] = q
		}
	}
	
	if exchangeName != "" {	
		// If the exchange name is undefined, use the default exchange
		if err = channel.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil); err != nil {
			panic(err)
		}

		// Bind queues to new exchange
		for _, name := range queueNames {
			if err = channel.QueueBind(queues[name].Name, name, exchangeName, false, nil); err != nil {
				panic(err)
			}
		}
	}

	return &AsyncRabbitMQClient{
		Connection: conn,
		Channel:	channel,
		Exchange:	exchangeName,
		Queues:     queues,
	}, nil
}

func (c *AsyncRabbitMQClient) Shutdown() {
	c.Connection.Close()
	c.Channel.Close()
}

func ConstructRabbitMQUrl(username, password, hostname, port string) string {
	return "amqp://" + username + ":" + password + "@" + hostname + ":" + port + "/"
}

func (c *AsyncRabbitMQClient) SendMessage(body, queueName string, timeoutSeconds int) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(timeoutSeconds) * time.Second,
	)
	defer cancel()
	
	publishing := amqp.Publishing{
		ContentType:  "text/plain",
		Body:		  []byte(body),
	}
	return c.Channel.PublishWithContext(ctx, "", c.Queues[queueName].Name, false,  false, publishing)
}

func (c *AsyncRabbitMQClient) ConsumeMessages(queueName string, outChannel chan string, stopConsumingChannel chan bool) {
	if messages, err := c.Channel.Consume(c.Queues[queueName].Name, "", true, false, false, false, nil); err != nil {
		return
	} else {
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
}