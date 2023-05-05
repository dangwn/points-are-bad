package api

import (
	// "log"
	// "net/http"
	// "strconv"

	"github.com/gin-gonic/gin"
)

/*
 * Structs
 */

/*
 * Router Methods
 */
 
func (r Router) addMatchGroup(rg *gin.RouterGroup) {
    points := rg.Group("/match")

	points.GET("/", getMatchesWithoutGoals)
}

func getMatchesWithoutGoals(c *gin.Context) {
	
}

/*
 * Services
 */