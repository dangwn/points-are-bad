package apiRouter

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

var (
	FRONTEND_DOMAIN string = getEnv("FRONTEND_DOMAIN", "localhost")

	CSRF_SECRET_KEY []byte = []byte(getEnv("CSRF_SECRET_KEY", "CSRF-Key"))
)