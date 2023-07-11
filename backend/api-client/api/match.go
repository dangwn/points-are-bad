package api

import (
	"fmt"
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
 
// Match router group
func (r Router) addMatchGroup(rg *gin.RouterGroup) {
    matchGroup := rg.Group("/match")

	matchGroup.GET("/", getMatchesWithoutGoals)
	matchGroup.POST("/", createMatch)
	matchGroup.PUT("/", updateMatch)
	matchGroup.DELETE("/", deleteMatch)
	matchGroup.GET("/full/", getFullMatches)
}

/*
 * Match Creation Endpoint (Admin Only)
 * Schema: {match_date: Date, home: string, away: string}
 * Returns: MatchWithId
 * Inserts a new match into DB with NULL home and away goals
 */
func createMatch(c *gin.Context) {
	if !isCurrentUserAdmin(c) {
		abortRouterMethod(c, http.StatusUnauthorized, "Not admin user", "User cannot create match as they are not an admin user")
		return
	}

	var match MatchWithoutGoals
	if err := c.BindJSON(&match); err != nil {
		logMessage := "Could not bind incoming match data in createMatch: " + err.Error()
		abortRouterMethod(c, http.StatusBadRequest, "Could not retrieve data from request", logMessage)
		return
	}

	newMatch, err := insertMatchIntoDb(match)
	if err != nil {
		logMessage := "Could not insert new match into DB in createMatch: " + err.Error()
		abortRouterMethod(c, http.StatusBadRequest, "Could not create match", logMessage)
		return
	}

	c.JSON(http.StatusCreated, newMatch)
	Logger.Info(fmt.Sprint("New match with ID ", newMatch.MatchId, " created"))
}

/*
 * Match Deletion Endpoint (Admin Only)
 * Schema: {match_id: int}
 * Deletes a match from the DB using a given match Id
 */
func deleteMatch(c *gin.Context) {
	if !isCurrentUserAdmin(c) {
		abortRouterMethod(c, http.StatusUnauthorized, "Not admin user", "User cannot delete match as they are not an admin user")
		return
	}

	var matchIdOnly MatchIdOnly
	if err := c.BindJSON(&matchIdOnly); err != nil {
		logMessage := "Could not get match Id from query in deleteMatch: " + err.Error()
		abortRouterMethod(c, http.StatusBadRequest, "Could not delete match", logMessage)
		return
	}

	if deleted, err := deleteMatchById(matchIdOnly.MatchId); !deleted {
		logMessage := "Could not delete match from DB in deleteMatch: " + err.Error()
		abortRouterMethod(c, http.StatusBadRequest, "Could not delete match", logMessage)
		return
	}

	c.Status(http.StatusNoContent)
	Logger.Info(fmt.Sprint("Match ", matchIdOnly.MatchId, " deleted successfully"))
}

/*
 * Full Matches Endpoint (Admin Only)
 * Query Params: start_date=Date&end_date=Date
 * Returns: []Match
 * Retrieves all matches with all information, within a given date range
 */
func getFullMatches(c *gin.Context) {
	if !isCurrentUserAdmin(c) {
		abortRouterMethod(c, http.StatusUnauthorized, "Not admin user", "User cannot get full matches as they are not an admin user")
		return
	}

	var dateRange DateRange
	if err := c.BindQuery(&dateRange); err != nil {
		logMessage := "Could not bind query in getFullMatches: " + err.Error()
		abortRouterMethod(c, http.StatusBadRequest, "Match dates in incorrect format", logMessage)
		return
	}

	matches, err := getFullMatchesInDateRange(dateRange.StartDate, dateRange.EndDate)
	if err != nil {
		logMessage := "Could not get matches from DB in getFullMatches: " + err.Error()
		abortRouterMethod(c, http.StatusBadRequest, "Could not retrieve matches", logMessage)
		return
	}

	c.JSON(http.StatusOK, matches)
	if userId, err := getCurrentUser(c); err == nil {
		Logger.Info("Full matches successfully fetched for user " + userId)
	} else {
		Logger.Info("Full matches returned to unknown user")
	}
}

/*
 * Matches Without Goals Endpoint
 * Query Params start_date=Date&end_date=Date
 * Returns: []MatchWithoutGoals
 * Retrieves all matches' match_dates, home and away teams, within a given date range
 */
func getMatchesWithoutGoals(c *gin.Context) {
	var dateRange DateRange
	if err := c.BindQuery(&dateRange); err != nil {
		logMessage := "Could not bind query in getMatchesWithoutGoals: " + err.Error()
		abortRouterMethod(c, http.StatusBadRequest, "Match dates in incorrect format", logMessage)
		return
	}

	matches, err := getMatchesInDateRange(dateRange.StartDate, dateRange.EndDate)
	if err != nil {
		logMessage := "Could not get matches from DB in getMatchesWithoutGoals: " + err.Error()
		abortRouterMethod(c, http.StatusBadRequest, "Could not retrieve matches", logMessage)
		return
	}

	c.JSON(http.StatusOK, matches)
	Logger.Info("Matches (without goals) successfully retrieved")
}

/*
 * Update Match Endpoint
 * Schema: {match_id: int, match_date: Date, home: string, away: string, home_goals: int|null, away_goals: int|null}
 * Returns: Match
 * Updates a match's data in DB, using a given match Id
 */
func updateMatch(c *gin.Context) {
	if !isCurrentUserAdmin(c) {
		abortRouterMethod(c, http.StatusUnauthorized, "User is not admin user", "User cannot update match as they are not an admin user")
		return
	}

	var match Match
	if err := c.BindJSON(&match); err != nil {
		logMessage := "Could not bind JSON in updateMatch: " + err.Error()
		abortRouterMethod(c, http.StatusBadRequest, "Could not get match from request", logMessage)
		return
	}

	if newMatch, err := updateMatchById(match); err != nil {
		logMessage := "Caught error when updating match in updateMatch: " + err.Error()
		abortRouterMethod(c, http.StatusBadRequest, "Could not update match", logMessage)
	} else if !newMatch {
		logMessage := "No rows affected when updating match in updateMatch: " + err.Error()
		abortRouterMethod(c, http.StatusBadRequest, "Could not update match", logMessage)
	} else {
		c.JSON(http.StatusAccepted, match)
		Logger.Info(fmt.Sprint("Match with Id ", match.MatchId, " updated successfully"))
	}
}

/*
 * Services
 */

/*
 * Deletes a match with a given Id in the DB
 * Only returns true when the number of rows affected is equal to 1
 * Otherwise it will return false even if the error is nil
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

// Gets []Match from the DB where startDate <= match_date < endDate
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

// Gets []MatchWithoutGoals from the DB where startDate <= match_date < endDate
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

/*
 * Inserts a match into the DB
 * If any errors are caught, an empty MatchWithId is returned
 */
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

/*
 * Updates a match in the DB
 * The match Id is supplied in the request and is used to update a particular match
 * Only returns true when the number of rows affected is equal to 1
 * Otherwise it will return false even if the error is nil
 */
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
		if n, err := result.RowsAffected(); n == 1 {
			return true, err
		} else {
			return false, err
		}
	}
}