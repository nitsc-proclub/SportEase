package repository

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/saku0512/GYOUJI_HP/backend/internal/model"
)

type TournamentRepository struct {
	DB *sql.DB
}

func NewTournamentRepository(db *sql.DB) *TournamentRepository {
	return &TournamentRepository{DB: db}
}

// GetTournamentsBySport fetches tournament data for a specific sport and converts it
// into the format required by the Bracketry library.
func (r *TournamentRepository) GetTournamentsBySport(sport string) (interface{}, error) {
	tournamentsRows, err := r.DB.Query("SELECT id, name FROM tournaments WHERE sport = ?", sport)
	if err != nil {
		return nil, err
	}
	defer tournamentsRows.Close()

	var finalBrackets []model.BracketryData

	for tournamentsRows.Next() {
		var tournamentID int64
		var tournamentName string
		if err := tournamentsRows.Scan(&tournamentID, &tournamentName); err != nil {
			return nil, err
		}

		teams, err := r.getTeamsForTournament(tournamentID)
		if err != nil {
			return nil, err
		}
		matches, err := r.getMatchesForTournament(tournamentID)
		if err != nil {
			return nil, err
		}

		bracket, err := r.transformToBracketryData(tournamentName, teams, matches)
		if err != nil {
			return nil, err
		}
		if bracket != nil {
			finalBrackets = append(finalBrackets, *bracket)
		}
	}

	if err := tournamentsRows.Err(); err != nil {
		return nil, err
	}

	return finalBrackets, nil
}

// transformToBracketryData is the core logic for converting DB data into the Bracketry format.
func (r *TournamentRepository) transformToBracketryData(tournamentName string, teams []model.Team, matches []model.Match) (*model.BracketryData, error) {
	if len(matches) == 0 {
		return nil, nil // 表示する試合がないので、このトーナメントはスキップ
	}
	matchesByRound := make(map[int][]model.Match)
	maxRound := 0
	for _, m := range matches {
		matchesByRound[m.Round] = append(matchesByRound[m.Round], m)
		if m.Round > maxRound {
			maxRound = m.Round
		}
	}

	var bracketryRounds []model.BracketryRound
	for i := 1; i <= maxRound; i++ {
		var roundName string
		// (ラウンド名のロジックはご提示いただいたものをそのまま使用します)
		if maxRound == 4 {
			switch i {
			case 1:
				roundName = "1st Round"
			case 2:
				roundName = "2nd Round"
			case 3:
				roundName = "Semifinals"
			case 4:
				roundName = "Final"
			}
		} else if maxRound == 3 {
			switch i {
			case 1:
				roundName = "1st Round"
			case 2:
				roundName = "Semifinals"
			case 3:
				roundName = "Final"
			}
		} else {
			roundName = fmt.Sprintf("Round %d", i)
		}
		bracketryRounds = append(bracketryRounds, model.BracketryRound{Name: roundName})
	}

	// Contestants を試合に参加するチームIDからも構築する（敗者戦など他トーナメント所属のチームIDにも対応）
	bracketryContestants := make(map[string]model.BracketryContestant)

	// まずはトーナメント所属チームを登録
	if len(teams) > 0 {
		for _, team := range teams {
			teamIDStr := strconv.FormatInt(team.ID, 10)
			bracketryContestants[teamIDStr] = model.BracketryContestant{
				EntryStatus: team.EntryStatus,
				Players:     []model.Player{{Title: team.Name}},
			}
		}
	}

	// 次に、試合に現れる全チームIDを収集
	teamIDSet := make(map[int64]struct{})
	for _, m := range matches {
		if m.Team1ID != nil {
			teamIDSet[*m.Team1ID] = struct{}{}
		}
		if m.Team2ID != nil {
			teamIDSet[*m.Team2ID] = struct{}{}
		}
		if m.WinnerTeamID != nil {
			teamIDSet[*m.WinnerTeamID] = struct{}{}
		}
	}
	// 既に登録済みのIDを除外
	var missingIDs []int64
	for id := range teamIDSet {
		idStr := strconv.FormatInt(id, 10)
		if _, exists := bracketryContestants[idStr]; !exists {
			missingIDs = append(missingIDs, id)
		}
	}
	// 未登録IDのチーム情報をID指定で取得し、Contestantsに追加
	if len(missingIDs) > 0 {
		fetched, err := r.getTeamsByIDs(missingIDs)
		if err != nil {
			return nil, err
		}
		for _, team := range fetched {
			teamIDStr := strconv.FormatInt(team.ID, 10)
			bracketryContestants[teamIDStr] = model.BracketryContestant{
				EntryStatus: team.EntryStatus,
				Players:     []model.Player{{Title: team.Name}},
			}
		}
	}

	var bracketryMatches []model.BracketryMatch
	for roundNum, roundMatches := range matchesByRound {
		sort.Slice(roundMatches, func(i, j int) bool {
			return roundMatches[i].MatchNumber < roundMatches[j].MatchNumber
		})

		// 3位決定戦が存在するラウンドかどうかのフラグ
		isFinalRoundWithBronze := (roundNum == maxRound && len(roundMatches) > 1)

		for order, match := range roundMatches {
			var sides []model.BracketrySide
			// (Side作成ロジックは変更ありません)
			if match.Team1ID != nil && match.Team2ID == nil {
				isWinner := true
				side1 := createSide(match.Team1ID, match.Team1Score, match.Team1ID)
				side1.IsWinner = &isWinner
				byeTitle := "<div style='margin-left: 60px'>BYE</div>"
				side2 := model.BracketrySide{Title: &byeTitle}
				sides = append(sides, side1, side2)
			} else if match.Team1ID == nil && match.Team2ID != nil {
				byeTitle := "<div style='margin-left: 60px'>BYE</div>"
				side1 := model.BracketrySide{Title: &byeTitle}
				isWinner := true
				side2 := createSide(match.Team2ID, match.Team2Score, match.Team2ID)
				side2.IsWinner = &isWinner
				sides = append(sides, side1, side2)
			} else {
				side1 := createSide(match.Team1ID, match.Team1Score, match.WinnerTeamID)
				side2 := createSide(match.Team2ID, match.Team2Score, match.WinnerTeamID)
				sides = append(sides, side1, side2)
			}

			// ★★★ 修正点2: 3位決定戦のロジックをここに追加 ★★★
			isBronze := false
			matchOrder := order // デフォルトの順序

			if isFinalRoundWithBronze {
				// ソート後の最初の試合（match_numberが小さい方）を3位決定戦とする
				if order == 0 {
					isBronze = true
					matchOrder = 1 // 3位決定戦のorderは1に固定
				} else { // 2番目の試合（match_numberが大きい方）を決勝戦とする
					isBronze = false
					matchOrder = 0 // 決勝戦のorderは0に固定
				}
			}

			bracketryMatches = append(bracketryMatches, model.BracketryMatch{
				RoundIndex:  roundNum - 1,
				Order:       matchOrder, // 調整後のorderを使用
				Sides:       sides,
				MatchStatus: match.Status,
				// ★★★ 修正点3: IsBronzeMatchフィールドを設定 ★★★
				IsBronzeMatch: isBronze,
			})
		}
	}

	return &model.BracketryData{
		Name:        tournamentName,
		Rounds:      bracketryRounds,
		Matches:     bracketryMatches,
		Contestants: bracketryContestants,
	}, nil
}

// createSide is a helper function to generate a side for a match.
func createSide(teamID *int64, score *int, winnerTeamID *int64) model.BracketrySide {
	side := model.BracketrySide{}

	if teamID != nil {
		idStr := strconv.FormatInt(*teamID, 10)
		side.ContestantID = &idStr
	}

	// The database has a single score, but the target format requires an array.
	// We wrap the single score in a slice.
	if score != nil {
		isScoreWinner := false
		if winnerTeamID != nil && teamID != nil && *winnerTeamID == *teamID {
			isScoreWinner = true
		}

		side.Scores = append(side.Scores, model.Score{
			MainScore: strconv.Itoa(*score),
			// NOTE: The 'isWinner' inside a score object is an assumption based on the target format.
			// You may need to adjust this logic based on your application's rules.
			IsWinner: &isScoreWinner,
		})
	}

	// Set the winner status for the entire side
	if winnerTeamID != nil && teamID != nil && *winnerTeamID == *teamID {
		isWinner := true
		side.IsWinner = &isWinner
	}
	return side
}

// getTeamsForTournament fetches all teams for a given tournament ID.
// NOTE: Assumes 'entry_status' and 'nationality' columns exist in the 'teams' table.
func (r *TournamentRepository) getTeamsForTournament(tournamentID int64) ([]model.Team, error) {
	rows, err := r.DB.Query("SELECT id, name, tournament_id, entry_status FROM teams WHERE tournament_id = ?", tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []model.Team
	for rows.Next() {
		var team model.Team
		if err := rows.Scan(&team.ID, &team.Name, &team.TournamentID, &team.EntryStatus); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	return teams, nil
}

// getMatchesForTournament fetches all matches for a given tournament ID.
// NOTE: Assumes a 'status' column exists in the 'matches' table.
func (r *TournamentRepository) getMatchesForTournament(tournamentID int64) ([]model.Match, error) {
	query := `
        SELECT id, tournament_id, round, match_number_in_round,
               team1_id, team2_id, team1_score, team2_score, winner_team_id, next_match_id, status
        FROM matches
        WHERE tournament_id = ?`
	rows, err := r.DB.Query(query, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []model.Match
	for rows.Next() {
		var match model.Match
		if err := rows.Scan(
			&match.ID, &match.TournamentID, &match.Round, &match.MatchNumber,
			&match.Team1ID, &match.Team2ID, &match.Team1Score, &match.Team2Score,
			&match.WinnerTeamID, &match.NextMatchID, &match.Status,
		); err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}
	return matches, nil
}

// getTeamsByIDs fetches team rows for the provided team IDs.
func (r *TournamentRepository) getTeamsByIDs(teamIDs []int64) ([]model.Team, error) {
	if len(teamIDs) == 0 {
		return []model.Team{}, nil
	}
	// 動的 IN 句を作成
	placeholders := make([]string, 0, len(teamIDs))
	args := make([]interface{}, 0, len(teamIDs))
	for _, id := range teamIDs {
		placeholders = append(placeholders, "?")
		args = append(args, id)
	}
	query := fmt.Sprintf("SELECT id, name, tournament_id, entry_status FROM teams WHERE id IN (%s)", strings.Join(placeholders, ","))
	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []model.Team
	for rows.Next() {
		var t model.Team
		if err := rows.Scan(&t.ID, &t.Name, &t.TournamentID, &t.EntryStatus); err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, nil
}
