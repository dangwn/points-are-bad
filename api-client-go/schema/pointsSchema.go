package schema

type Score struct {
	Points        int  `json:"points"`
	CorrectScores int  `json:"correct_scores"`
	LargestError  int  `json:"largest_error"`
	Position      *int `json:"position"`
}

type UserWithPoints struct {
	Score
	SessionUser
}

type LeaderBoardUser struct {
	Score
	Username string `json:"username"`
}