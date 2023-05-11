package api

import (
	"log"
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

	q := NewQuery(
		"matches",
		"prediction_id",
		"users_and_predictions.home_goals",
		"users_and_predictions.away_goals",
		"match_date",
		"home",
		"away",
		"matches.home_goals",
		"matches.away_goals",
		"username",
		"is_admin",
	).Join(
		NewQuery(
			"predictions", "prediction_id", "home_goals", "away_goals", "username", "is_admin", "match_id",
		).Join(
			"users", "inner", "predictions.user_id", "users.user_id",
		).Filter(
			"users.user_id", "=", userId,
		).NameQuery(
			"users_and_predictions",
		), 
		"inner", 
		"matches.match_id", 
		"users_and_predictions.match_id",
	)
	
	if startDate != nil { q = q.Filter("match_date", ">=", startDate) }
	if endDate != nil { q = q.Filter("match_date", "<", endDate) }

	rows, err := q.All()

	if err != nil{
		return predictions, err
	}

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

	return predictions, nil
}