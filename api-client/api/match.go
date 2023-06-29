package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
 * Structs
 */
type Match struct {
	MatchId string `json:"match_id"`
	MatchWithoutGoals
	HomeGoals *int `json:"home_goals"`
	AwayGoals *int `json:"away_goals"`
}

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

type MatchIdOnly struct {
	MatchId int `json:"match_id" form:"match_id"`
}

/*
 * Router Methods
 */
 
func (r Router) addMatchGroup(rg *gin.RouterGroup) {
    matchGroup := rg.Group("/match")

	matchGroup.GET("/", getMatchesWithoutGoals)
	matchGroup.POST("/", createMatch)
	matchGroup.PUT("/", updateMatch)
	matchGroup.DELETE("/", deleteMatch)
	matchGroup.GET("/full/", getFullMatches)
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

func deleteMatch(c *gin.Context) {
	if !isCurrentUserAdmin(c) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"detail":"Not admin user",
		})
		return
	}

	var matchIdOnly MatchIdOnly
	if err := c.BindQuery(&matchIdOnly); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Could not get match ID from request",
		})
		return
	}

	if deleted, _ := deleteMatchById(matchIdOnly.MatchId); !deleted {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Could not delete match",
		})
		return
	}
	c.Status(http.StatusNoContent)
}

func getFullMatches(c *gin.Context) {
	if !isCurrentUserAdmin(c) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"detail":"Not admin user",
		})
		return
	}

	var dateRange DateRange
	if err := c.BindQuery(&dateRange); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Incorrect date formats",
		})
		return
	}

	matches, err := getFullMatchesInDateRange(dateRange.StartDate, dateRange.EndDate)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Could not retrieve matches",
		})
		return
	}

	c.JSON(http.StatusAccepted, matches)
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
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Could not retrieve matches",
		})
	}

	c.JSON(http.StatusAccepted, matches)
}

func updateMatch(c *gin.Context) {
	if !isCurrentUserAdmin(c) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"detail":"Not admin user",
		})
		return
	}

	var match Match
	if err := c.BindJSON(&match); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Could not get match ID from request",
		})
		return
	}

	if newMatch, err := updateMatchById(match); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Could not update match",
		})
	} else if !newMatch {
		log.Println("could not update match")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Could not find match to update",
		})
	} else {
		c.JSON(http.StatusAccepted, match)
	}
}

/*
 * Services
 */

func deleteMatchById(matchId int) (bool, error) {
	deleteQuery := "DELETE FROM matches WHERE match_id = $1"
	if result, err := driver.Exec(deleteQuery, matchId); err != nil {
		return false, err
	} else {
		if n, err := result.RowsAffected(); n == 1 {
			return true, err 
		} else {
			return false, err
		}
	}
}

func getFullMatchesInDateRange(startDate *Date, endDate *Date) ([]Match, error) {
	var matches []Match

	matchQuery := "SELECT * FROM matches" + createDateRangeWhereClause(startDate, endDate) + " ORDER BY match_date"

	if rows, err := driver.Query(matchQuery); err != nil {
		return matches, err
	} else {
		for rows.Next() {
			var match Match
			if err := rows.Scan(
				&match.MatchId,
				&match.MatchDate,
				&match.Home,
				&match.Away,
				&match.HomeGoals,
				&match.AwayGoals,
			); err != nil {
				return matches, err
			}
			matches = append(matches, match)
		}
	}

	return matches, nil
}

func getMatchesInDateRange(startDate *Date, endDate *Date) ([]MatchWithoutGoals, error) {
	var matches []MatchWithoutGoals

	matchQuery := "SELECT match_date, home, away FROM matches" + createDateRangeWhereClause(startDate, endDate) + " ORDER BY match_date"

	if rows, err := driver.Query(matchQuery); err != nil {
		return matches, err
	} else {
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
	}

	return matches, nil	
}

func insertMatchIntoDb(match MatchWithoutGoals) (MatchWithId, error) {
	var matchId int

	insertQuery := `
		INSERT INTO matches(match_date, home, away) 
		VALUES ($1, $2, $3) 
		RETURNING match_id
	`
	if err := driver.QueryRow(
		insertQuery, 
		match.MatchDate, 
		match.Home, 
		match.Away,
	).Scan(&matchId); err != nil {
		return MatchWithId{}, err
	}

	if err := populatePredictionsByMatchId(matchId); err != nil {
		return MatchWithId{}, err
	}

	return MatchWithId{
		matchId,
		match,
	}, nil
}

func updateMatchById(match Match) (bool, error) {
	if result, err := driver.Exec(
		`UPDATE matches 
		SET 
			home = $1,
			away = $2,
			match_date = $3,
			home_goals = $4,
			away_goals = $5
		WHERE match_id = $6`,
		match.Home,
		match.Away,
		match.MatchDate,
		match.HomeGoals,
		match.AwayGoals,
		match.MatchId,
	); err != nil {
		return false, err
	} else {
		if affected, err := result.RowsAffected(); err != nil {
			return false, err
		} else {
			if affected > 0 {return true, nil}
			return false, nil
		}
	}
}