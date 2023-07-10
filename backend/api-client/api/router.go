package api

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	router *gin.Engine
}

func NewRouter() Router {
	r := Router{
		router: gin.Default(),
	}
	r.router.Use(corsMiddleware())

	baseGroup := r.router.Group("")
	r.addUserGroup(baseGroup)
	r.addAuthGroup(baseGroup)
	r.addPointsGroup(baseGroup)
	r.addPredictionGroup(baseGroup)
	r.addMatchGroup(baseGroup)

    r.router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message":"let's go"})
    })

	return r
}

func (r Router) Run(addr ...string) error {
	return r.router.Run(addr...)
}
