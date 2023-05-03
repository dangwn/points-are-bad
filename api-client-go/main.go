package main

import (
	"points-areb-bad/api-client/api-router"
)

func main() {
	router := apiRouter.NewRouter()
	router.Run(":8000")
}
