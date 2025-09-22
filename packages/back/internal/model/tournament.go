package model

// このファイルは、フロントエンドのライブラリ 'bracketry' が要求するデータ構造を定義します。

type BracketryData struct {
	Name        string                         `json:"name"`
	Rounds      []BracketryRound               `json:"rounds"`
	Matches     []BracketryMatch               `json:"matches"`
	Contestants map[string]BracketryContestant `json:"contestants"`
}

// BracketryRound defines the name of a round.
type BracketryRound struct {
	Name string `json:"name"`
}

// BracketryMatch represents a single match.
type BracketryMatch struct {
	RoundIndex    int             `json:"roundIndex"`
	Order         int             `json:"order"`
	Sides         []BracketrySide `json:"sides"`
	MatchStatus   *string         `json:"matchStatus,omitempty"`
	IsBronzeMatch bool            `json:"isBronzeMatch"`
}

// BracketrySide represents one side of a match (a player or a BYE).
type BracketrySide struct {
	ContestantID *string `json:"contestantId,omitempty"`
	Scores       []Score `json:"scores,omitempty"`
	IsWinner     *bool   `json:"isWinner,omitempty"`
	Title        *string `json:"title,omitempty"` // Used for BYEs
}

// Score represents the score within a set or game. The target format requires string.
type Score struct {
	MainScore string `json:"mainScore"`
	Subscore  *int   `json:"subscore,omitempty"`
	IsWinner  *bool  `json:"isWinner,omitempty"`
}

type TeamScore struct {
	Team1Score *int `json:"team1_score" binding:"required"`
	Team2Score *int `json:"team2_score" binding:"required"`
}

// BracketryContestant represents a participant (a team or player).
type BracketryContestant struct {
	EntryStatus *string  `json:"entryStatus,omitempty"`
	Players     []Player `json:"players"`
}

// Player represents a single player within a contestant entry.
type Player struct {
	Title string `json:"title"`
}

// --- データベースから読み込むための構造体 ---

// Team represents the data fetched from the 'teams' table.
type Team struct {
	ID           int64
	Name         string
	TournamentID int64
	EntryStatus  *string // Assumed new column in 'teams' table
}

// Match represents the data fetched from the 'matches' table.
type Match struct {
	ID           int64   `json:"id"`
	TournamentID int64   `json:"tournament_id"`
	Round        int     `json:"round"`
	MatchNumber  int     `json:"match_number_in_round"` // Corresponds to match_number_in_round
	Team1ID      *int64  `json:"team1_id"`
	Team2ID      *int64  `json:"team2_id"`
	Team1Score   *int    `json:"team1_score"`
	Team2Score   *int    `json:"team2_score"`
	WinnerTeamID *int64  `json:"winner_team_id"`
	NextMatchID  *int64  `json:"next_match_id"`
	Status       *string `json:"status"` // Assumed new column in 'matches' table
}

type ErrorResponse struct {
	Error string `json:"error" example:"具体的なエラーメッセージ"`
}
