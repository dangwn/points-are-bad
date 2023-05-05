package main

import (
	"os"

	apiRouter "points-are-bad/api-client/api"
)

var API_PORT string = getEnv("API_PORT", "8020")

func getEnv(varName string, defaultValue string) string {
	envVar := os.Getenv(varName)
	if envVar == "" {
		return defaultValue
	}	
	return envVar
}

func main() {
	router := apiRouter.NewRouter()
	router.Run(":"+API_PORT)
}
