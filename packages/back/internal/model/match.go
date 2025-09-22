package model

// MatchResponse は /api/matches/{sport} のレスポンスで使われる試合情報の構造体です。
// チーム名など、追加情報を含みます。
type MatchResponse struct {
	ID                 int64   `json:"id"`
	MatchNumberInRound int64   `json:"match_number_in_round"`
	Round              int64   `json:"round"`
	Team1ID            *int64  `json:"team1_id"`
	Team1Name          *string `json:"team1_name"`
	Team2ID            *int64  `json:"team2_id"`
	Team2Name          *string `json:"team2_name"`
	Team1Score         *int    `json:"team1_score"`
	Team2Score         *int    `json:"team2_score"`
	WinnerTeamID       *int64  `json:"winner_team_id"`
	Status             *string `json:"status"`
	NextMatchID        *int64  `json:"next_match_id"`
	TournamentName     *string `json:"tournament_name,omitempty"`
}

type SuccessResponse struct {
	Message string `json:"message" example:"処理が正常に完了しました"`
}
