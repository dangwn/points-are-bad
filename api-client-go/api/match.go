package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
 * Structs
 */
type MatchWithoutGoals struct {
	MatchDate *Date   `json:"match_date"`
	Home 	  string  `json:"home"`
	Away 	  string  `json:"away"`
}

type MatchWithId struct {
	MatchId int `json:"match_id"`
	MatchWithoutGoals
}

type DateRange struct {
	StartDate *Date `json:"start_date" query:"start_date" form:"start_date"`
	EndDate *Date `json:"end_date" query:"end_date" form:"end_date"`
}

/*
 * Router Methods
 */
 
func (r Router) addMatchGroup(rg *gin.RouterGroup) {
    matchGroup := rg.Group("/match")

	matchGroup.GET("/", getMatchesWithoutGoals)
	matchGroup.POST("/", createMatch)
}

func createMatch(c *gin.Context) {
	if !isCurrentUserAdmin(c) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"detail":"Not admin user",
		})
		return
	}

	var match MatchWithoutGoals
	if err := c.BindJSON(&match); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Could not retrieve data from request",
		})
		return
	}

	if newMatch, err := insertMatchIntoDb(match); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Could not create match",
		})
	} else {
		c.JSON(http.StatusCreated, newMatch)
	}
}

func getMatchesWithoutGoals(c *gin.Context) {
	var dateRange DateRange
	if err := c.BindQuery(&dateRange); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Incorrect date formats",
		})
		return
	}

	matches, err := getMatchesInDateRange(dateRange.StartDate, dateRange.EndDate)
	log.Println(matches)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Could not retrieve matches",
		})
	}

	c.JSON(http.StatusAccepted, matches)
}

/*
 * Services
 */

func insertMatchIntoDb(match MatchWithoutGoals) (MatchWithId, error) {
	var matchId int

	if err := driver.InsertWithReturn(
		"matches",
		"match_date, home, away",
		"$1, $2, $3",
		"match_id",
		match.MatchDate,
		match.Home,
		match.Away,
	).Scan(&matchId); err != nil {
		return MatchWithId{}, err
	}

	return MatchWithId{
		matchId,
		match,
	}, nil
}

func getMatchesInDateRange(startDate *Date, endDate *Date) ([]MatchWithoutGoals, error) {
	var matches []MatchWithoutGoals
	var rows *sql.Rows
	var err error
	baseQuery := "SELECT match_date, home, away FROM matches "

	if startDate != nil {
		if endDate != nil {
			rows, err = driver.Query(
				baseQuery + "WHERE match_date >= $1 AND match_date < $2",
				startDate,
				endDate,
			)
		} else {
			rows, err = driver.Query(
				baseQuery + "WHERE match_date >= $1",
				startDate,
			)
		}
	} else {
		if endDate != nil {
			rows, err = driver.Query(
				baseQuery + "WHERE match_date < $1",
				endDate,
			)
		} else {
			rows, err = driver.Query(baseQuery)
		}
	}
	if err != nil {
		return matches, err
	}

	for rows.Next() {
		var match MatchWithoutGoals
		if err := rows.Scan(
			&match.MatchDate,
			&match.Home,
			&match.Away,
		); err != nil {
			return matches, err
		}

		matches = append(matches, match)
	}

	return matches, nil
}