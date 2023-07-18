package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

/*
* Structs
 */

type Score struct {
    Points          int  `json:"points"`
    CorrectScores   int  `json:"correct_scores"`
    LargestError    int  `json:"largest_error"`
    Position        *int `json:"position"`
}

type LeaderboardUser struct {
    Score
    Username string `json:"username"`
}

type LimitOffset struct {
    Limit   int `json:"limit" query:"limit" form:"limit"`
    Offset  int `json:"offset" query:"offset" form:"offset"`
}

/*
 * Router Methods
 */

// Points router group
func (r Router) addPointsGroup(rg *gin.RouterGroup) {
    points := rg.Group("/points")

    points.GET("/", getUserPoints)
    points.GET("/leaderboard/", getLeaderboard)
    points.POST("/calculate/", calculatePoints)
}

/*
 * Points Endpoint
 * Returns: Score
 * Retrieves the user's points
 */
func getUserPoints(c *gin.Context) {
    currentUserId, err := getCurrentUser(c)
    if err != nil {
        logMessage := "Could not retrieve current user in getUserPoints: " + err.Error()
        abortRouterMethod(c, http.StatusUnauthorized, "Could not retrieve current user", logMessage)
        return
    }

    userPoints, err := getPointsByUserId(currentUserId)
    if err != nil {
        logMessage := "Could not get user's points in getUserPoints: " + err.Error()
        abortRouterMethod(c, http.StatusUnauthorized, "Could not retrieve user's points", logMessage)
        return
    }

    c.JSON(http.StatusOK, userPoints)
    Logger.Info("User " + currentUserId + " successfully requested points")
} 

/*
 * Leaderboard Endpoint
 * Query Params: limit=int&offset=int
 * Returns: []LeaderboardUser
 * Retrieves the leaderboard with a given number of entries and offset
 */
func getLeaderboard(c *gin.Context) {
    var lo LimitOffset
    if err := c.BindQuery(&lo); err != nil {
        logMessage := "Could not bind query in getLeaderBoard: " + err.Error()
        abortRouterMethod(c, http.StatusBadRequest, "Could not retrieve leaderboard", logMessage)
        return
    }

    leaderboard, err := getGlobalLeaderboard(lo.Limit, lo.Offset)
    if err != nil {
        logMessage := "Could not retrieve leaderboard in getLeaderBoard: " + err.Error()
        abortRouterMethod(c, http.StatusBadRequest, "Could not retrieve leaderboard", logMessage)
        return
    }

    c.JSON(http.StatusOK, leaderboard)
    Logger.Info("Leaderboard successfully retrieved")
}

/*
 * Points Calculation Endpoint (Admin Only)
 * Calculates all player's points and updates DB with new points/correct scores/largest error/position
 */
func calculatePoints(c *gin.Context) {
    if !isCurrentUserAdmin(c) {
        abortRouterMethod(c, http.StatusUnauthorized, "Not admin user", "User cannot calculate match as they are not an admin user")
        return
    }

    if err := updatePoints(); err != nil {
        logMessage := "Could not update points in calculatePoints: " + err.Error()
        abortRouterMethod(c, http.StatusUnauthorized, "Could not update points", logMessage)
        return
    }

    c.Status(http.StatusNoContent)
    if userId, err := getCurrentUser(c); err == nil {
        Logger.Info("User " + userId + " updated global points")
    } else {
        Logger.Info("Unkown user updated global points")
    }
}

/*
 * Services
 */

// Runs SQL query to get leaderboard from "offset" with "limit" entries
func getGlobalLeaderboard(limit, offset int) ([]LeaderboardUser, error) {
    var leaderboard []LeaderboardUser

    pointsQuery := `
        SELECT points, correct_scores, largest_error, position, username
        FROM users
        ORDER BY position
        LIMIT $1 
        OFFSET $2
    `
    if rows, err := driver.Query(pointsQuery, limit, offset); err != nil {
        return leaderboard, err
    } else {
        for rows.Next() {
            var leaderboardRow LeaderboardUser
            if err := rows.Scan(
                &leaderboardRow.Points,
                &leaderboardRow.CorrectScores,
                &leaderboardRow.LargestError,
                &leaderboardRow.Position,
                &leaderboardRow.Username,
            ); err != nil {
                return leaderboard, err
            }
            leaderboard = append(leaderboard, leaderboardRow)
        }
    }

    return leaderboard, nil
}

// Retrieves a user's points using given user Id
func getPointsByUserId(userId string) (Score, error) {
    var userPoints Score

    pointsQuery := `
        SELECT points, correct_scores, largest_error, position
        FROM users
        WHERE users.user_id = $1
        LIMIT 1
    `
    err := driver.QueryRow(pointsQuery, userId).Scan(
        &userPoints.Points,
        &userPoints.CorrectScores,
        &userPoints.LargestError,
        &userPoints.Position,
    )

    return userPoints, err
}

/*
 * Runs query to update all users' points, correct scores, largest error and position
 * Applies a given penalty, NULL_PREDICTIONS_PENALTY (see config.go), when users have not predicted a score
 * Runs the query up to today's date
 * Will not calculate points when a match has NULL home_goals or away_goals
 * NOTE: The update can still be thrown even if the query has successfully completed, when
 *         the result.RowsAffected() produces an error
 */
func updatePoints() error {
    todayDate := time.Now().Format("2001-01-01")

    updatePointsQuery := `UPDATE users
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
                                SELECT user_id, home_goals as pred_hg, away_goals as pred_ag, match_id
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
        WHERE users.user_id = new_points.user_id`

    result, err := driver.Exec(updatePointsQuery, NULL_PREDICTIONS_PENALTY, todayDate)
    if err != nil {
        return err
    }
    if _ ,err := result.RowsAffected(); err != nil {
        return err
    }

    return nil
}