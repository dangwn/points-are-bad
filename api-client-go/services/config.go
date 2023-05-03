package services

import (
	"os"
	"time"
)

func getEnv(varName string, defaultValue string) string {
	envVar := os.Getenv(varName)
	if envVar == "" {
		return defaultValue
	}	
	return envVar
}

var (
	ACCESS_TOKEN_SECRET_KEY []byte = []byte(getEnv(
		"ACCESS_TOKEN_SECRET_KEY", 
		"OOOOO-ACCESS-SECRET",
	))
	ACCESS_TOKEN_EXPIRE_TIME time.Duration = time.Minute * 15

	REFRESH_TOKEN_SECRET_KEY []byte = []byte(getEnv(
		"REFRESH_TOKEN_SECRET_KEY", 
		"OOOOO-REFRESH-SECRET",
	))
	REFRESH_TOKEN_MAX_AGE_DAYS = 30
	REFRESH_TOKEN_EXPIRE_TIME time.Duration = time.Hour * 24 * time.Duration(REFRESH_TOKEN_MAX_AGE_DAYS)

	BCRYPT_SECRET_KEY []byte = []byte(getEnv(
		"BCRYPT_SECRET_KEY",
		"OOOOO-BCRYPT-SECRET",
	))

	POSTGRES_USER string = getEnv("POSTGRES_USER", "test")
	POSTGRES_PASSWORD string = getEnv("POSTGRES_PASSWORD", "password")
	POSTGRES_DB string = getEnv("POSTGRES_DB", "db")

	RABBITMQ_USER string = getEnv("RABBITMQ_USER", "guest")
	RABBITMQ_PASSWORD string = getEnv("RABBITMQ_PASSWORD", "guest")
	RABBITMQ_HOST string = getEnv("RABBITMQ_HOST", "localhost")
	RABBITMQ_PORT string = getEnv("RABBITMQ_PORT", "5672")

	EMAIL_VERIFICATION_QUEUE string = "email-verification"

	LETTERS []rune = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	REDIS_HOST string = getEnv("REDIS_HOST", "localhost")
	REDIS_PORT string = getEnv("REDIS_PORT", "6379")
	REDIS_PASSWORD string = getEnv("REDIS_PASSWORD", "")
	REDIS_DB int = 0
)