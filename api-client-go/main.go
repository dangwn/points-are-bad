package main

import (
	"points-are-bad/api-client/api-router"
)

func main() {
	router := apiRouter.NewRouter()
	router.Run(":8000")
}
