package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS middleware
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Request origin: ", c.Request.Header.Get("Origin"))

		// For swagger docs to send requests in local build
		if IS_DEV_BUILD {
			fmt.Println("Dev build - allowing all origins")
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", FRONTEND_DOMAIN)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}