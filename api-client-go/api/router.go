package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	router *gin.Engine
}

func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", FRONTEND_DOMAIN)
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
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
	r.addMatchGroup(baseGroup)

    r.router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message":"let's go"})
    })

	return r
}

func (r Router) Run(addr ...string) error {
	return r.router.Run(addr...)
}
