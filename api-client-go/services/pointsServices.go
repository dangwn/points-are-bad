package services

import (
	"points-are-bad/api-client/schema"
)

func GetGlobalLeaderboard(limit int, offset int) ([]schema.LeaderBoardUser, error) {
	var leaderboard []schema.LeaderBoardUser

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
		var leaderboardRow schema.LeaderBoardUser
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

	return leaderboard, nil
}

func GetPointsByUserId(userId string) (schema.UserWithPoints, error) {
	var userPoints schema.UserWithPoints
	
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
		&userPoints.Username,
		&userPoints.IsAdmin,
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
