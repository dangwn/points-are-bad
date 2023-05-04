package apiRouter

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	// "points-are-bad/api-client/schema"
	"points-are-bad/api-client/services"
)

func (r Router) addPointsGroup(rg *gin.RouterGroup) {
    points := rg.Group("/points")

	points.GET("/", getUserPoints)
	points.GET("/leaderboard", getGlobalLeaderboard)
}

func getUserPoints(c *gin.Context) {
	currentUserId, err := services.GetCurrentUser(c)
    if err != nil {
        log.Println(err)
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "detail":"Could not retreieve current user",
        })
        return
    }

	userPoints, err := services.GetPointsByUserId(currentUserId)
	if err != nil {
        log.Println(err)
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "detail":"Could not retreieve current user's points",
        })
        return
    }

	c.JSON(http.StatusOK, userPoints)
} 

func getGlobalLeaderboard(c *gin.Context) {
	offsetString, limitString := c.Query("offset"), c.Query("limit")
	if offsetString == "" {
		offsetString = "0"
	}
	if limitString == "" {
		limitString = "10"
	}

	limit, err := strconv.Atoi(limitString)
	if err != nil {
		log.Println(err)
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "detail":"Limit was not provided as an integer",
        })
	}
	offset, err := strconv.Atoi(offsetString)
	if err != nil {
		log.Println(err)
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "detail":"Offset was not provided as an integer",
        })
	}

	leaderboard, err := services.GetGlobalLeaderboard(limit, offset)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "detail":"Could not retrieve leaderboard",
        })
	}

	c.JSON(http.StatusOK, leaderboard)
}