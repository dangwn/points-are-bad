package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
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

func (r Router) addPredictionGroup(rg *gin.RouterGroup) {
	prediction := rg.Group("/prediction")

	prediction.GET("/", getUserPredictions)
}

func getUserPredictions(c *gin.Context) {
	currentUserId, err := getCurrentUser(c)
    if err != nil {
        log.Println(err)
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "detail":"Could not retreieve current user",
        })
        return
    }

	if userPredictions, err := getPredictionsByUserId(currentUserId, nil, nil); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "detail":"Could not retreieve predictions",
        })
	} else {
		c.JSON(http.StatusOK, userPredictions)
	}
}

/*
 *  Services
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

func populatePredictionsByUserId(userId string) error {
	var matchIds []int

	if rows, err := driver.Query("SELECT match_id FROM matches"); err != nil {
		return err
	} else {
		for rows.Next() {
			var matchId int
			if err := rows.Scan(&matchId); err != nil {
				return err
			}

			matchIds = append(matchIds, matchId)
		}
	}	
	if len(matchIds) == 0 {
		return nil
	}

	insertQuery := `
		INSERT INTO predictions(home_goals, away_goals, user_id, match_id)
		VALUES (NULL, NULL, $1, $2)
	`
	_, err := driver.Exec(insertQuery, userId, pq.Array(matchIds))
	return err
}

func populatePredictionsByMatchId(matchId int) error {
	var userIds []string

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

	insertQuery := `
		INSERT INTO predictions(home_goals, away_goals, user_id, match_id)
		VALUES (NULL, NULL, $1, $2)
	`
	_, err := driver.Exec(insertQuery, pq.Array(userIds), matchId)
	return err
}