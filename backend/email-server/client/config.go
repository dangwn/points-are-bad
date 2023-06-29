package emailClient

import (
	"os"
)

func getEnv(varName string, defaultValue string) string {
	envVar := os.Getenv(varName)
	if envVar == "" {
		return defaultValue
	}	
	return envVar
}

func createSignUpTokenUrl() string {
	host := getEnv("WEB_HOST_NAME", "http://localhost:3000")
	if host[len(host) - 1] != '/' {
		host += "/"
	}
	return host + "signup?token="
}

var (
	AWS_ACCESS_KEY_ID string = os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY string = os.Getenv("AWS_SECRET_ACCESS_KEY")
	AWS_REGION string = getEnv("AWS_REGION", "us-east-1")
	AWS_SESSION_TOKEN string = getEnv("AWS_SESSION_TOKEN", "")

	EMAIL_VERIFICATION_QUEUE_NAME string = "email-verification"
	EMAIL_BUFFER_SIZE int = 1000
	
	SENDER_EMAIL string = getEnv("SENDER_EMAIL", "")

	SIGNUP_TOKEN_URL string = createSignUpTokenUrl()

	RABBITMQ_USER string = getEnv("RABBITMQ_USER", "guest")
	RABBITMQ_PASSWORD string = getEnv("RABBITMQ_PASSWORD", "guest")
	RABBITMQ_HOST string = getEnv("RABBITMQ_HOST", "localhost")
	RABBITMQ_PORT string = getEnv("RABBITMQ_PORT", "5672")
)