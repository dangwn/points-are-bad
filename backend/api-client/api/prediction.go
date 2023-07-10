package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
 *  Structs
 */
type Prediction struct {
	PredictionId int `json:"prediction_id"`
	HomeGoals *int `json:"home_goals"`
	AwayGoals *int `json:"away_goals"`
}

type PredictionWithMatch struct {
	Prediction
	PredictionMatch Match `json:"match"`
	PredictionUser SessionUser `json:"user"`
}


/*
 *  Router Methods
 */

// Prediction router group
func (r Router) addPredictionGroup(rg *gin.RouterGroup) {
	prediction := rg.Group("/prediction")

	prediction.GET("/", getUserPredictions)
	prediction.PUT("/", updateUserPredictions)
}

/*
 * User Predictions Predictions
 * Returns: []Prediction
 * Retrieves all of a user's predictions
 */
func getUserPredictions(c *gin.Context) {
	currentUserId, err := getCurrentUser(c)
    if err != nil {
        logMessage := "Could not retrieve current user in getUserPredictions: " + err.Error()
		abortRouterMethod(c, http.StatusUnauthorized, "Could not get user predictions", logMessage)
		return
    }

	userPredictions, err := getPredictionsByUserId(currentUserId, nil, nil)
	if err != nil {
		logMessage := "Could not retrieve user predictions in getUserPredicitons: " + err.Error()
		abortRouterMethod(c, http.StatusBadRequest, "Could not get user predictions", logMessage)
		return
	}

	c.JSON(http.StatusOK, userPredictions)
	Logger.Info("User " + currentUserId + " successfully requested predictions")
}

/*
 * Update Predictions Endpoint
 * Schema: [{prediction_id: int, home_goals: int, away_goals: int}, ...]
 * Updates a user's predictions with new supplied predictions
 */
func updateUserPredictions(c *gin.Context) {
	currentUserId, err := getCurrentUser(c)
    if err != nil {
        logMessage := "Could not retrieve current user in updateUserPredictions: " + err.Error()
		abortRouterMethod(c, http.StatusUnauthorized, "Could not get user predictions", logMessage)
		return
    }

	predictions := new([]Prediction)
	if body, err := io.ReadAll(c.Request.Body); err != nil {
		logMessage := "Could not read request body in updateUserPredictions: " + err.Error()
		abortRouterMethod(c, http.StatusBadRequest, "Could not get user predictions", logMessage)
		return
	} else {
		if err := json.Unmarshal(body, &predictions); err != nil {
			logMessage := "Could not unmarshal request body in updateUserPredictions: " + err.Error()
			abortRouterMethod(c, http.StatusBadRequest, "Could not get user predictions", logMessage)
			return
		}
	}

	if err := updatePredictionsByUserId(*predictions, currentUserId); err != nil {
		logMessage := "Could not update predictions in updateUserPredictions: " + err.Error()
		abortRouterMethod(c, http.StatusBadRequest, "Could not get user predictions", logMessage)
		return
	}

	c.Status(http.StatusAccepted)
	Logger.Info("User " + currentUserId + " updated their predictions")
}

/*
 *  Services
 */

/*
 * Retrieves user predictions from DB in a given date range
 * If no dates are supplied, all predictions are returned
 */
func getPredictionsByUserId(userId string, startDate *Date, endDate *Date) ([]PredictionWithMatch, error) {
	var predictions []PredictionWithMatch

	predictionsQuery := `
		WITH uap AS (
			SELECT prediction_id, home_goals, away_goals, username, is_admin, match_id
			FROM predictions 
			INNER JOIN users 
				ON predictions.user_id = users.user_id
			WHERE predictions.user_id = $1
		)
		SELECT prediction_id, uap.home_goals, uap.away_goals, match_date,
			home, away, matches.home_goals, matches.away_goals, username, is_admin
		FROM matches
		INNER JOIN uap 
			ON matches.match_id = uap.match_id
	` + createDateRangeWhereClause(startDate, endDate)

	if rows, err := driver.Query(predictionsQuery, userId); err != nil {
		return predictions, err
	} else {
		for rows.Next() {
			var pred PredictionWithMatch
			if err := rows.Scan(
				&pred.PredictionId,
				&pred.HomeGoals,
				&pred.AwayGoals,
				&pred.PredictionMatch.MatchDate,
				&pred.PredictionMatch.Home,
				&pred.PredictionMatch.Away,
				&pred.PredictionMatch.HomeGoals,
				&pred.PredictionMatch.AwayGoals,
				&pred.PredictionUser.Username,
				&pred.PredictionUser.IsAdmin,
			); err != nil {
				return predictions, err
			}
			predictions = append(predictions, pred)
		}
	}

	return predictions, nil
}

/*
 * Populates predictions table in DB with entries for a given user Id
 * This is used when a new user is created, and they will have prediction entries
 * 		created for every match currently in the DB
 */
func populatePredictionsByUserId(userId string) error {
	var matchIds []int64

	// Get match Ids for matches currently in DB
	if rows, err := driver.Query("SELECT match_id FROM matches"); err != nil {
		return err
	} else {
		for rows.Next() {
			var matchId int64
			if err := rows.Scan(&matchId); err != nil {
				return err
			}
			matchIds = append(matchIds, matchId)
		}
	}	
	if len(matchIds) == 0 {
		return nil
	}

	// Unnests match Ids into query string (see customTypes.go)
	insertQuery := `
		INSERT INTO predictions(home_goals, away_goals, user_id, match_id)
		SELECT NULL, NULL, $1, ` + UnnestArray(matchIds).String()
		                         
	_, err := driver.Exec(insertQuery, userId)
	return err
}

/*
 * Populates predictions table in DB with entries for a given match Id
 * This is used when a new match is created, and all users will have a new prediction
 * 		entry created for the new match
 */
func populatePredictionsByMatchId(matchId int) error {
	var userIds []string

	// Get user Ids for matches currently in DB
	if rows, err := driver.Query("SELECT user_id FROM users"); err != nil {
		return err
	} else {
		for rows.Next() {
			var userId string
			if err := rows.Scan(&userId); err != nil {
				return err
			}

			userIds = append(userIds, userId)
		}
	}	
	if len(userIds) == 0 {
		return nil
	}

	// Unnests user Ids into query string (see customTypes.go)
	insertQuery := `
		INSERT INTO predictions(home_goals, away_goals, match_id, user_id)
		SELECT NULL, NULL, $1, ` + UnnestArray(userIds).String()

	_, err := driver.Exec(insertQuery, matchId)
	return err
}

/*
 * Updates a user's predictions with an array of new predictions
 * Verifies that all supplied predictions have that user's user Id
 * 		If the verification fails, an error is thrown
 * NOTE: Because a check is done to see whether all predictions have been changed,
 *			some predictions may be changed and an error is still thrown
 */
func updatePredictionsByUserId(predictions []Prediction, userId string) error {
	// Start constructing query string in goroutine
	queryStringChannel := make(chan string)
	go func() {
		queryStringChannel <- `
		UPDATE predictions AS p SET
			prediction_id = c.prediction_id,
			home_goals = c.home_goals,
			away_goals = c.away_goals
		FROM (
			VALUES ` + PredictionArray(predictions).String() + `
		) AS c(prediction_id, home_goals, away_goals)
		WHERE 
			p.prediction_id = c.prediction_id
		`
	}()
	
	// Check that all predictions in "predictions" belong to the user
	// The predicitons Ids are stored in the keys of a map, with an empty struct as their value
	userPredictionIds := make(map[int]struct{})
	if rows, err := driver.Query("SELECT prediction_id FROM predictions WHERE user_id = $1", userId); err != nil {
		return err
	} else {
		for rows.Next() {
			var predId int
			if err := rows.Scan(&predId); err != nil {
				return err
			}
			userPredictionIds[predId] = struct{}{}
		}
	}
	// Iterate through predictions and verify a key with that Id exists in the map
	// If even one mismatch is found, throw an error
	for _, v := range predictions {
		if _, found := userPredictionIds[v.PredictionId]; !found {
			return errors.New("prediction being changed that does not belong to user " + userId)
		}
	}

	// Retrieve query string from channel
	queryString := <- queryStringChannel

	// Execute query and check that all predictions have been changed
	if result, err := driver.Exec(queryString); err != nil {
		return err
	} else {
		if affected, err := result.RowsAffected(); err != nil {
			return err
		} else if affected != int64(len(predictions)) {
			return errors.New("different number of rows affected than expected")
		}
	}

	return nil
}