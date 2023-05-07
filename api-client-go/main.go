package main

import (
	"points-are-bad/api-client/api"
)

func main() {
	router := api.NewRouter()
	router.Run(":" + api.API_PORT)
}
