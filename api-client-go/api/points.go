package api

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

/*
* Structs
 */

type Score struct {
	Points        int  `json:"points"`
	CorrectScores int  `json:"correct_scores"`
	LargestError  int  `json:"largest_error"`
	Position      *int `json:"position"`
}

type UserWithPoints struct {
	Score
	User SessionUser `json:"user"`
}

type LeaderBoardUser struct {
	Score
	User Username `json:"user"`
}

/*
 * Router Methods
 */

func (r Router) addPointsGroup(rg *gin.RouterGroup) {
    points := rg.Group("/points")

	points.GET("/", getUserPoints)
	points.GET("/leaderboard/", getLeaderboard)
	points.POST("/calculate/", calculatePoints)
}

func getUserPoints(c *gin.Context) {
	currentUserId, err := getCurrentUser(c)
    if err != nil {
        log.Println(err)
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "detail":"Could not retreieve current user",
        })
        return
    }

	userPoints, err := getPointsByUserId(currentUserId)
	if err != nil {
        log.Println(err)
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "detail":"Could not retreieve current user's points",
        })
        return
    }

	c.JSON(http.StatusOK, userPoints)
} 

func getLeaderboard(c *gin.Context) {
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

	leaderboard, err := getGlobalLeaderboard(limit, offset)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "detail":"Could not retrieve leaderboard",
        })
	}

	c.JSON(http.StatusOK, leaderboard)
}

func calculatePoints(c *gin.Context) {
	if isAdmin := isCurrentUserAdmin(c); !isAdmin {
		log.Println("user is not admin")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err := updatePoints(); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Could not update points",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

/*
 * Services
 */
func getGlobalLeaderboard(limit int, offset int) ([]LeaderBoardUser, error) {
	var leaderboard []LeaderBoardUser

	rows, err := driver.Query(
		`SELECT points, correct_scores, largest_error, position, username
		FROM points 
		LEFT JOIN users 
		ON users.user_id = points.user_id
		ORDER BY position
		LIMIT $1
		OFFSET $2`,
		limit,
		offset,
	)
	if err != nil {
		return leaderboard, err
	}

	for rows.Next() {
		var leaderboardRow LeaderBoardUser
		if err := rows.Scan(
			&leaderboardRow.Points,
			&leaderboardRow.CorrectScores,
			&leaderboardRow.LargestError,
			&leaderboardRow.Position,
			&leaderboardRow.User.Username,
		); err != nil {
			return leaderboard, err
		}
		leaderboard = append(leaderboard, leaderboardRow)
	}

	return leaderboard, nil
}

func getPointsByUserId(userId string) (UserWithPoints, error) {
	var userPoints UserWithPoints
	
	err := driver.QueryRow(
		`SELECT points, correct_scores, largest_error, position, username, is_admin
		FROM points 
		LEFT JOIN users 
		ON users.user_id = points.user_id
		WHERE points.user_id = $1`,
		userId,
	).Scan(
		&userPoints.Points,
		&userPoints.CorrectScores,
		&userPoints.LargestError,
		&userPoints.Position,
		&userPoints.User.Username,
		&userPoints.User.IsAdmin,
	) 

	return userPoints, err
}

func insertPointsIntoDb(userId string) error {
	_, err := driver.Insert(
		"points",
		"user_id, points, correct_scores, largest_error, position",
		"$1, 0, 0, 0, NULL",
		userId,
	)
	return err
}

func updatePoints() error {
	todayDate := time.Now().Format("2001-01-01")

	updatePointsQuery := `UPDATE points
		SET points = new_points.points, correct_scores = new_points.correct_scores, largest_error = new_points.largest_error, position = new_points.position
		FROM (
				SELECT *, RANK() OVER (ORDER BY points ASC, correct_scores DESC, largest_error ASC) as position
				FROM (    
					SELECT user_id, SUM(points) as points, COUNT(CASE WHEN points = 0 THEN 1 END) as correct_scores, MAX(points) as largest_error
					FROM (
						SELECT user_id, COALESCE(ABS(pred_hg-hg) + ABS(pred_ag-ag), $1) as points
						FROM (
							SELECT user_id, pred_hg, pred_ag, home_goals as hg, away_goals as ag
							FROM (  
								SELECT "user_id", "home_goals" as "pred_hg", "away_goals" as "pred_ag", match_id
								FROM predictions
							) as t1
							JOIN matches on t1.match_id = matches.match_id
							WHERE 
								matches.match_date < $2 AND
								matches.home_goals IS NOT NULL AND
								matches.away_goals IS NOT NULL
						) as t2
					) as t3
					GROUP BY user_id
				) as t4
		) as new_points
		WHERE points.user_id = new_points.user_id`

	result, err := driver.Exec(updatePointsQuery, NULL_PREDICTIONS_PENALTY, todayDate)
	if err != nil {
		return err
	}
	if _ ,err := result.RowsAffected(); err != nil {
		return err
	}

	return nil
}