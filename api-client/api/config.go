package api

import (
	"os"
	"strings"
	"time"
)

func getEnv(varName string, defaultValue string) string {
	envVar := os.Getenv(varName)
	if envVar == "" {
		return defaultValue
	}	
	return envVar
}

func isProductionBuild() bool {
	envName := os.Getenv("PAB_ENVIRONMENT")
	return strings.ToLower(envName) == "production"
}

var (
	ACCESS_TOKEN_SECRET_KEY []byte = []byte(getEnv("ACCESS_TOKEN_SECRET_KEY", "OOOOO-ACCESS-SECRET"))
	ACCESS_TOKEN_EXPIRE_TIME time.Duration = time.Minute * 15
	
	API_PORT string = getEnv("API_PORT", "8020")

	BCRYPT_SECRET_KEY []byte = []byte(getEnv(
		"BCRYPT_SECRET_KEY",
		"OOOOO-BCRYPT-SECRET",
	))

	FRONTEND_DOMAIN string = getEnv("FRONTEND_DOMAIN", "http://localhost:3000")

	EMAIL_VERIFICATION_QUEUE string = "email-verification"

	IS_DEV_BUILD bool = !isProductionBuild()

	CHARS []rune = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	NULL_PREDICTIONS_PENALTY int = 10

	POSTGRES_USER string = getEnv("POSTGRES_USER", "test")
	POSTGRES_PASSWORD string = getEnv("POSTGRES_PASSWORD", "password")
	POSTGRES_DB string = getEnv("POSTGRES_DB", "db")
	
	RABBITMQ_USER string = getEnv("RABBITMQ_USER", "guest")
	RABBITMQ_PASSWORD string = getEnv("RABBITMQ_PASSWORD", "guest")
	RABBITMQ_HOST string = getEnv("RABBITMQ_HOST", "localhost")
	RABBITMQ_PORT string = getEnv("RABBITMQ_PORT", "5672")

	REDIS_HOST string = getEnv("REDIS_HOST", "localhost")
	REDIS_PORT string = getEnv("REDIS_PORT", "6379")
	REDIS_PASSWORD string = getEnv("REDIS_PASSWORD", "")
	REDIS_DB int = 0
	REDIS_DURATION time.Duration = time.Minute * 30

	REFRESH_TOKEN_NAME string = "X-Refresh-Token"
	REFRESH_TOKEN_SECRET_KEY []byte = []byte(getEnv("REFRESH_TOKEN_SECRET_KEY", "OOOOO-REFRESH-SECRET"))
	REFRESH_TOKEN_MAX_AGE_DAYS int = 30
	REFRESH_TOKEN_EXPIRE_TIME time.Duration = time.Hour * 24 * time.Duration(REFRESH_TOKEN_MAX_AGE_DAYS)
	REFRESH_TOKEN_EXPIRE_TIME_SECONDS int = (3600 * 24) * REFRESH_TOKEN_MAX_AGE_DAYS
)