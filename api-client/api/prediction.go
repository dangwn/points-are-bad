package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

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

type PredictionArray []Prediction

func (p PredictionArray) String() string {
	if n := len(p); n == 0 {
		return ""
	} else if n > 0 {
		// Value takes form "(id, h, a),(id, h, a),..."
		b := make([]byte, 0, 8*n-2)
		b = appendPredictionToArrayBuffer(b, p[0])
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = appendPredictionToArrayBuffer(b, p[i])
		}		

		return string(b)
	}

	return ""
}


/*
 *  Router Methods
 */

func (r Router) addPredictionGroup(rg *gin.RouterGroup) {
	prediction := rg.Group("/prediction")

	prediction.GET("/", getUserPredictions)
	prediction.PUT("/", updateUserPredictions)
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

func updateUserPredictions(c *gin.Context) {
	currentUserId, err := getCurrentUser(c)
    if err != nil {
        log.Println(err)
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "detail":"Could not retreieve current user",
        })
        return
    }

	body, _ := io.ReadAll(c.Request.Body)
	predictions := new([]Prediction)
	if err := json.Unmarshal(body, &predictions); err != nil {
		log.Println(err)
		return
	}

	if err := updatePredictionsByUserId(*predictions, currentUserId); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"detail":"Could not update user predictions",
		})
		return
	}

	c.Status(http.StatusAccepted)
}

/*
 *  Services
 */
func appendIntPointerToSqlBuffer(b []byte, i *int) []byte {
	if i == nil {
		return append(b, 'N', 'U', 'L', 'L')
	}
	return strconv.AppendInt(b, int64(*i), 10)
}

func appendPredictionToArrayBuffer(b []byte, pred Prediction) []byte {
	b = append(b, '(')
	b = strconv.AppendInt(b, int64(pred.PredictionId), 10)
	b = append(b, ',')
	b = appendIntPointerToSqlBuffer(b, pred.HomeGoals)
	b = append(b, ',')
	b = appendIntPointerToSqlBuffer(b, pred.AwayGoals)
	return append(b, ')')
}

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
	var matchIds []int64

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

	insertQuery := `
		INSERT INTO predictions(home_goals, away_goals, user_id, match_id)
		SELECT NULL, NULL, $1, `+ UnnestArray(matchIds).String()
		                         
	_, err := driver.Exec(insertQuery, userId)
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
		INSERT INTO predictions(home_goals, away_goals, match_id, user_id)
		SELECT NULL, NULL, $1, ` + UnnestArray(userIds).String()

	_, err := driver.Exec(insertQuery, matchId)
	return err
}

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
	for _, v := range predictions {
		if _, found := userPredictionIds[v.PredictionId]; !found {
			return errors.New("prediction being changed that does not belong to user")
		}
	}

	queryString := <- queryStringChannel

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