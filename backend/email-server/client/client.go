package emailClient

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	rmq "github.com/dangwn/points-are-bad-tooling/rmqclient"
)

type EmailServer struct {
	Service *ses.SES
	Rabbit 	*rmq.AsyncRabbitMQClient
}

type VerificationMessage struct {
	Email 			  string `json:"email"` 
	VerificationToken string `json:"token"`
}

// Creates a new email server - consisting of an SES Email Server and a RabbitMQ client
func NewEmailServer(accessKey, secretAccessKey, region, token string, rabbitmqClient *rmq.AsyncRabbitMQClient) (*EmailServer, error) {
	if sess, err := session.NewSession(
		&aws.Config{
			Credentials: credentials.NewStaticCredentials(
				accessKey, secretAccessKey, token,
			),
			Region: aws.String(region),
		},
	); err != nil {
		return nil, err
	} else {
		return &EmailServer{
			Service: ses.New(sess),
			Rabbit: rabbitmqClient,
		}, nil
	}
}

// Sends a given email using AWS SES
func (server *EmailServer) SendEmail(body, html, emailSubject, recipientEmail, senderEmail string) error {
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(recipientEmail)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data: aws.String(html),
				},
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data: aws.String(body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data: aws.String(emailSubject),
			},
		},
		Source: aws.String(senderEmail),
	}
	_, err := server.Service.SendEmail(input)
	return err
}

/*
 * Function to be used as a goroutine to handle verification email requests
 * 		coming from rabbitmq
 * "inChannel" is set to be the output channel from the rabbitmq listener
 * Decodes message, expecting {"email":..., "token":...}
 */
func (server *EmailServer) HandleIncomingVerificationRequests(inChannel chan string) {
	Logger.Info("Email handling running")
	for {
		msg := <- inChannel
		data := VerificationMessage{}
		json.Unmarshal([]byte(msg), &data)

		if (data.Email == "" || data.VerificationToken == "") {
			Logger.Warning("Could not unpack message from RabbitMQ")
			continue
		}

		body := createEmailBody(data.VerificationToken)
		html := createEmailHtml(data.VerificationToken)

		Logger.Info("Sending verification email")
		if err := server.SendEmail(body, html, EMAIL_SUBJECT, data.Email, SENDER_EMAIL); err != nil {
			Logger.Warning(err)
		}
	}
}

// Function to be used as goroutine to listen for messages coming from rabbitmq
func (server *EmailServer) ListenForVerificationRequests(queueName string, outChannel chan string, stopConsumingChannel chan bool) error {
	Logger.Info("Listening on queue \"" + queueName + "\"")
	return server.Rabbit.ConsumeMessages(queueName, outChannel, stopConsumingChannel)
}

// Runs the email server and starts listening for messages from RabbitMQ
func (server *EmailServer) Run(queueName string) {
	listenOutChannel := make(chan string, EMAIL_BUFFER_SIZE)
	listenStopChannel := make(chan bool)

	go server.ListenForVerificationRequests(queueName, listenOutChannel, listenStopChannel)
	go server.HandleIncomingVerificationRequests(listenOutChannel)
	Logger.Info("Server running...")

	<- listenStopChannel
}