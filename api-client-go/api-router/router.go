package apiRouter 

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
)

type Router struct {
	router *gin.Engine
}

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost, http://localhost:3000")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

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
	r.router.Use(CORSMiddleware())

	baseGroup := r.router.Group("")
	r.addUserGroup(baseGroup)
	r.addAuthGroup(baseGroup)

	return r
}

func (r Router) Run(addr ...string) error {
	return r.router.Run(addr...)
}