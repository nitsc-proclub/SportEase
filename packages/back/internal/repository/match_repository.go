package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/saku0512/GYOUJI_HP/backend/internal/model"
)

type MatchRepository struct {
	DB *sql.DB
}

func NewMatchRepository(db *sql.DB) *MatchRepository {
	return &MatchRepository{DB: db}
}

// UpdateMatchScore は試合のスコアを更新し、勝者を次の試合へ進出させます。
func (r *MatchRepository) UpdateMatchScore(matchID int, team1Score int, team2Score int, user interface{}) (interface{}, error) {
	userMap, ok := user.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	role, _ := userMap["role"].(string)
	assignedSport, _ := userMap["assigned_sport"].(string)

	tx, err := r.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			log.Printf("Rolling back transaction due to error: %v", err)
			tx.Rollback()
		}
	}()

	var tournamentID int
	var team1ID, team2ID, nextMatchID, loserNextMatchID, previousWinnerID sql.NullInt64
	var matchSport, tournamentName string
	var currentRound int
	query := `
        SELECT m.tournament_id, m.round, t.name, m.team1_id, m.team2_id, 
               m.next_match_id, m.loser_next_match_id, t.sport, m.winner_team_id
        FROM matches m JOIN tournaments t ON m.tournament_id = t.id
        WHERE m.id = ?`
	err = tx.QueryRow(query, matchID).Scan(&tournamentID, &currentRound, &tournamentName, &team1ID, &team2ID, &nextMatchID, &loserNextMatchID, &matchSport, &previousWinnerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("not found")
		}
		return nil, fmt.Errorf("failed to query initial match: %w", err)
	}

	if !(role == "superroot" || (role == "admin" && assignedSport == matchSport)) {
		return nil, fmt.Errorf("forbidden")
	}

	if previousWinnerID.Valid {
		if role != "superroot" {
			return nil, fmt.Errorf("match result already submitted")
		}
		if err := r.resetMatchResult(tx, matchID); err != nil {
			return nil, fmt.Errorf("failed to reset match result: %w", err)
		}
	}

	var winnerTeamID, loserTeamID sql.NullInt64
	var winnerClassID, loserClassID sql.NullInt64
	if team1Score > team2Score {
		winnerTeamID = team1ID
		loserTeamID = team2ID
	} else if team2Score > team1Score {
		winnerTeamID = team2ID
		loserTeamID = team1ID
	}
	// チームIDからclass_idを取得
	if winnerTeamID.Valid {
		err = tx.QueryRow("SELECT class_id FROM teams WHERE id = ?", winnerTeamID.Int64).Scan(&winnerClassID)
		if err != nil {
			// エラー処理: クエリが失敗したか、結果がなかった場合
			log.Printf("Could not find class_id for winner team %d: %v", winnerTeamID.Int64, err)
		}
	}
	if loserTeamID.Valid {
		err = tx.QueryRow("SELECT class_id FROM teams WHERE id = ?", loserTeamID.Int64).Scan(&loserClassID)
		if err != nil {
			log.Printf("Could not find class_id for loser team %d: %v", loserTeamID.Int64, err)
		}
	}

	// DRY原則に基づき、更新ロジックを関数化
	err = r.updateMatchAndProgress(tx, matchID, team1Score, team2Score, winnerTeamID, loserTeamID, nextMatchID, loserNextMatchID)
	if err != nil {
		return nil, err
	}

	// --- ポイント付与ロジック ---
	isLoserBracket := tournamentName == "卓球（雨天時・敗者戦左側）" || tournamentName == "卓球（雨天時・敗者戦右側）"

	if winnerTeamID.Valid && winnerClassID.Valid && !isLoserBracket {
		// ラウンドごとに該当カラムに10点を登録
		var scoreColumn string
		switch currentRound {
		case 1:
			scoreColumn = matchSport + "1_score"
		case 2:
			scoreColumn = matchSport + "2_score"
		case 3:
			scoreColumn = matchSport + "3_score"
		}
		if scoreColumn != "" {
			query := fmt.Sprintf("UPDATE team_points SET %s = 10 WHERE class_id = ?", scoreColumn)
			_, err := tx.Exec(query, winnerClassID.Int64)
			if err != nil {
				return nil, err
			}
		}

		// 決勝戦・3位決定戦のみchampionship_score加算
		var maxRound int
		_ = tx.QueryRow(`SELECT MAX(round) FROM matches WHERE tournament_id = ?`, tournamentID).Scan(&maxRound)

		if currentRound == maxRound {
			var maxNum, thisNum int
			_ = tx.QueryRow(`SELECT MAX(match_number_in_round) FROM matches WHERE tournament_id = ? AND round = ?`, tournamentID, currentRound).Scan(&maxNum)
			_ = tx.QueryRow(`SELECT match_number_in_round FROM matches WHERE id = ?`, matchID).Scan(&thisNum)
			champCol := matchSport + "_championship_score"
			// 決勝戦（ラウンド最終試合）のみ加算
			if thisNum == maxNum {
				// 決勝戦
				var championshipScore int = 80
				query := fmt.Sprintf("UPDATE team_points SET %s = ? WHERE class_id = ?", champCol)
				_, err := tx.Exec(query, championshipScore, winnerClassID.Int64)
				if err != nil {
					return nil, err
				}
				if loserTeamID.Valid && loserClassID.Valid {
					championshipScore = 60
					_, err := tx.Exec(query, championshipScore, loserClassID.Int64)
					if err != nil {
						return nil, err
					}
				}
			} else if thisNum == maxNum-1 && maxNum > 1 {
				// 3位決定戦（max-1番目）
				var bronzeScore int = 50
				query := fmt.Sprintf("UPDATE team_points SET %s = ? WHERE class_id = ?", champCol)
				_, err := tx.Exec(query, bronzeScore, winnerClassID.Int64)
				if err != nil {
					return nil, err
				}
				if loserTeamID.Valid && loserClassID.Valid {
					bronzeScore = 40
					_, err := tx.Exec(query, bronzeScore, loserClassID.Int64)
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}

	// 卓球（晴天時）の場合、雨天時も同期更新するロジック
	if matchSport == "table_tennis" && tournamentName == "卓球（晴天時）" {
		var rainyMatchID int
		var rainyNextMatchID, rainyLoserNextMatchID sql.NullInt64
		var rainyTeam1ID, rainyTeam2ID sql.NullInt64

		rainyQuery := `
			SELECT m_rainy.id, m_rainy.next_match_id, m_rainy.loser_next_match_id, m_rainy.team1_id, m_rainy.team2_id, m_rainy.match_number_in_round
			FROM matches m_sunny
			JOIN tournaments t_rainy ON t_rainy.name = '卓球（雨天時）'
			JOIN matches m_rainy ON m_rainy.tournament_id = t_rainy.id 
				AND m_rainy.round = m_sunny.round 
				AND m_rainy.match_number_in_round = m_sunny.match_number_in_round
			WHERE m_sunny.id = ?`

		var rainyMatchNumberInRound sql.NullInt64
		err = tx.QueryRow(rainyQuery, matchID).Scan(&rainyMatchID, &rainyNextMatchID, &rainyLoserNextMatchID, &rainyTeam1ID, &rainyTeam2ID, &rainyMatchNumberInRound)
		if err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("failed to find corresponding rainy match: %w", err)
		}

		if rainyMatchID != 0 {
			if previousWinnerID.Valid { // If we are editing a result
				if err := r.resetMatchResult(tx, rainyMatchID); err != nil {
					return nil, fmt.Errorf("failed to reset rainy match result: %w", err)
				}
			}

			// 雨天側の勝者/敗者は雨天トーナメントのチームIDで決定する
			var rainyWinnerTeamID, rainyLoserTeamID sql.NullInt64
			if team1Score > team2Score {
				rainyWinnerTeamID = rainyTeam1ID
				rainyLoserTeamID = rainyTeam2ID
			} else if team2Score > team1Score {
				rainyWinnerTeamID = rainyTeam2ID
				rainyLoserTeamID = rainyTeam1ID
			}

			err = r.updateMatchAndProgress(tx, rainyMatchID, team1Score, team2Score, rainyWinnerTeamID, rainyLoserTeamID, rainyNextMatchID, rainyLoserNextMatchID)
			if err != nil {
				return nil, fmt.Errorf("failed to sync rainy match: %w", err)
			}

			// 1回戦の敗者を雨天の敗者復活戦トーナメントへ登録（重複防止チェックつき）
			if currentRound == 1 && rainyLoserTeamID.Valid && rainyMatchNumberInRound.Valid {
				var loserTournamentName string
				var loserEntryMatchNumber int
				switch rainyMatchNumberInRound.Int64 {
				case 1, 2:
					loserTournamentName = "卓球（雨天時・敗者戦左側）"
					loserEntryMatchNumber = 13
				case 3, 4:
					loserTournamentName = "卓球（雨天時・敗者戦左側）"
					loserEntryMatchNumber = 14
				case 5, 6:
					loserTournamentName = "卓球（雨天時・敗者戦右側）"
					loserEntryMatchNumber = 13
				case 7, 8:
					loserTournamentName = "卓球（雨天時・敗者戦右側）"
					loserEntryMatchNumber = 14
				}
				if loserTournamentName != "" {
					var loserEntryMatchID sql.NullInt64
					q := `
						SELECT m.id
						FROM tournaments tr
						JOIN matches m ON m.tournament_id = tr.id
						WHERE tr.name = ? AND m.round = 1 AND m.match_number_in_round = ?`
					err = tx.QueryRow(q, loserTournamentName, loserEntryMatchNumber).Scan(&loserEntryMatchID)
					if err == nil && loserEntryMatchID.Valid {
						// 既に同じチームが登録済みか確認し、未登録なら割り当て
						var existingTeam1, existingTeam2 sql.NullInt64
						if err := tx.QueryRow(`SELECT team1_id, team2_id FROM matches WHERE id = ?`, loserEntryMatchID.Int64).Scan(&existingTeam1, &existingTeam2); err == nil {
							if (existingTeam1.Valid && existingTeam1.Int64 == rainyLoserTeamID.Int64) || (existingTeam2.Valid && existingTeam2.Int64 == rainyLoserTeamID.Int64) {
								// すでに登録済み
							} else {
								if err := r.progressTeamToNextMatch(tx, rainyLoserTeamID, loserEntryMatchID); err != nil {
									return nil, fmt.Errorf("failed to progress loser to rainy consolation: %w", err)
								}
							}
						}
					}
				}
			}
		}
	}

	// 雨天時の敗者復活戦ブロックで1位になった場合のボーナス加算
	if winnerTeamID.Valid && winnerClassID.Valid && matchSport == "table_tennis" && (tournamentName == "卓球（雨天時・敗者戦左側）" || tournamentName == "卓球（雨天時・敗者戦右側）") {
		// 決勝ラウンド（1位決定戦）のみ加算
		var cnt int
		_ = tx.QueryRow(`SELECT COUNT(*) FROM matches WHERE tournament_id = ? AND round > ?`, tournamentID, currentRound).Scan(&cnt)
		if cnt == 0 {
			_, err := tx.Exec("UPDATE team_points SET table_tennis_rainy_bonus_score = 10 WHERE class_id = ?", winnerClassID.Int64)
			if err != nil {
				return nil, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return r.getMatchByID(matchID)
}

// 以前の勝者・敗者を次の試合から削除し、関連する得点をすべて0にリセットする
func (r *MatchRepository) resetMatchResult(tx *sql.Tx, matchID int) error {
	var tournamentID, currentRound int
	var team1ID, team2ID, nextMatchID, loserNextMatchID, previousWinnerID sql.NullInt64
	var matchSport, tournamentName string

	query := `
        SELECT m.tournament_id, m.round, t.name, m.team1_id, m.team2_id, 
               m.next_match_id, m.loser_next_match_id, t.sport, m.winner_team_id
        FROM matches m JOIN tournaments t ON m.tournament_id = t.id
        WHERE m.id = ?`
	err := tx.QueryRow(query, matchID).Scan(&tournamentID, &currentRound, &tournamentName, &team1ID, &team2ID, &nextMatchID, &loserNextMatchID, &matchSport, &previousWinnerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("match not found for reset: %d", matchID)
		}
		return fmt.Errorf("failed to query match for reset: %w", err)
	}

	if !previousWinnerID.Valid {
		return nil // Nothing to reset
	}

	var previousLoserID sql.NullInt64
	if team1ID.Valid && previousWinnerID.Int64 == team1ID.Int64 {
		previousLoserID = team2ID
	} else if team2ID.Valid && previousWinnerID.Int64 == team2ID.Int64 {
		previousLoserID = team1ID
	}

	if nextMatchID.Valid {
		if _, err := tx.Exec("UPDATE matches SET team1_id = NULL WHERE id = ? AND team1_id = ?", nextMatchID.Int64, previousWinnerID.Int64); err != nil {
			return err
		}
		if _, err := tx.Exec("UPDATE matches SET team2_id = NULL WHERE id = ? AND team2_id = ?", nextMatchID.Int64, previousWinnerID.Int64); err != nil {
			return err
		}
	}
	if previousLoserID.Valid && loserNextMatchID.Valid {
		if _, err := tx.Exec("UPDATE matches SET team1_id = NULL WHERE id = ? AND team1_id = ?", loserNextMatchID.Int64, previousLoserID.Int64); err != nil {
			return err
		}
		if _, err := tx.Exec("UPDATE matches SET team2_id = NULL WHERE id = ? AND team2_id = ?", loserNextMatchID.Int64, previousLoserID.Int64); err != nil {
			return err
		}
	}

	var previousWinnerClassID, previousLoserClassID sql.NullInt64
	if previousWinnerID.Valid {
		_ = tx.QueryRow("SELECT class_id FROM teams WHERE id = ?", previousWinnerID.Int64).Scan(&previousWinnerClassID)
	}
	if previousLoserID.Valid {
		_ = tx.QueryRow("SELECT class_id FROM teams WHERE id = ?", previousLoserID.Int64).Scan(&previousLoserClassID)
	}

	isLoserBracket := tournamentName == "卓球（雨天時・敗者戦左側）" || tournamentName == "卓球（雨天時・敗者戦右側）"

	if previousWinnerClassID.Valid && !isLoserBracket {
		var scoreColumn string
		switch currentRound {
		case 1:
			scoreColumn = matchSport + "1_score"
		case 2:
			scoreColumn = matchSport + "2_score"
		case 3:
			scoreColumn = matchSport + "3_score"
		}
		if scoreColumn != "" {
			query := fmt.Sprintf("UPDATE team_points SET %s = 0 WHERE class_id = ?", scoreColumn)
			if _, err := tx.Exec(query, previousWinnerClassID.Int64); err != nil {
				return err
			}
		}

		var maxRound int
		_ = tx.QueryRow(`SELECT MAX(round) FROM matches WHERE tournament_id = ?`, tournamentID).Scan(&maxRound)
		if currentRound == maxRound {
			var maxNum, thisNum int
			_ = tx.QueryRow(`SELECT MAX(match_number_in_round) FROM matches WHERE tournament_id = ? AND round = ?`, tournamentID, currentRound).Scan(&maxNum)
			_ = tx.QueryRow(`SELECT match_number_in_round FROM matches WHERE id = ?`, matchID).Scan(&thisNum)
			champCol := matchSport + "_championship_score"
			if thisNum == maxNum || (thisNum == maxNum-1 && maxNum > 1) {
				query := fmt.Sprintf("UPDATE team_points SET %s = 0 WHERE class_id = ?", champCol)
				if _, err := tx.Exec(query, previousWinnerClassID.Int64); err != nil {
					return err
				}
				if previousLoserClassID.Valid {
					if _, err := tx.Exec(query, previousLoserClassID.Int64); err != nil {
						return err
					}
				}
			}
		}
	}

	if previousWinnerClassID.Valid && matchSport == "table_tennis" && isLoserBracket {
		var cnt int
		_ = tx.QueryRow(`SELECT COUNT(*) FROM matches WHERE tournament_id = ? AND round > ?`, tournamentID, currentRound).Scan(&cnt)
		if cnt == 0 {
			if _, err := tx.Exec("UPDATE team_points SET table_tennis_rainy_bonus_score = 0 WHERE class_id = ?", previousWinnerClassID.Int64); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *MatchRepository) updateMatchAndProgress(tx *sql.Tx, matchID int, team1Score int, team2Score int, winnerTeamID, loserTeamID, nextMatchID, loserNextMatchID sql.NullInt64) error {
	_, err := tx.Exec(`UPDATE matches SET team1_score = ?, team2_score = ?, winner_team_id = ?, status = 'finished' WHERE id = ?`,
		team1Score, team2Score, winnerTeamID, matchID)
	if err != nil {
		return fmt.Errorf("failed to update match id %d: %w", matchID, err)
	}
	if winnerTeamID.Valid && nextMatchID.Valid {
		err := r.progressTeamToNextMatch(tx, winnerTeamID, nextMatchID)
		if err != nil {
			return fmt.Errorf("failed to progress winner from match id %d: %w", matchID, err)
		}
	}
	if loserTeamID.Valid && loserNextMatchID.Valid {
		err := r.progressTeamToNextMatch(tx, loserTeamID, loserNextMatchID)
		if err != nil {
			return fmt.Errorf("failed to progress loser from match id %d: %w", matchID, err)
		}
	}
	return nil
}

func (r *MatchRepository) progressTeamToNextMatch(tx *sql.Tx, teamToProgress, nextMatchID sql.NullInt64) error {
	var team1, team2 sql.NullInt64
	err := tx.QueryRow(`SELECT team1_id, team2_id FROM matches WHERE id = ?`, nextMatchID.Int64).Scan(&team1, &team2)
	if err != nil {
		return err
	}

	var updateQuery string
	if !team1.Valid {
		updateQuery = `UPDATE matches SET team1_id = ? WHERE id = ?`
	} else if !team2.Valid {
		updateQuery = `UPDATE matches SET team2_id = ? WHERE id = ?`
	}

	if updateQuery != "" {
		_, err = tx.Exec(updateQuery, teamToProgress.Int64, nextMatchID.Int64)
		if err != nil {
			return err
		}
	}
	return nil
}

// getMatchByID は指定されたIDの試合情報を取得するヘルパー関数です。
func (r *MatchRepository) getMatchByID(matchID int) (interface{}, error) {
	var match model.Match
	query := `
        SELECT id, tournament_id, round, match_number_in_round,
               team1_id, team2_id, team1_score, team2_score, winner_team_id, next_match_id
        FROM matches
        WHERE id = ?`
	err := r.DB.QueryRow(query, matchID).Scan(
		&match.ID, &match.TournamentID, &match.Round, &match.MatchNumber,
		&match.Team1ID, &match.Team2ID, &match.Team1Score, &match.Team2Score,
		&match.WinnerTeamID, &match.NextMatchID,
	)
	if err != nil {
		return nil, err
	}
	return match, nil
}

// GetMatchesBySport は指定された競技の試合一覧を取得します。
func (r *MatchRepository) GetMatchesBySport(sport string) ([]model.MatchResponse, error) {
	query := `
		SELECT m.id, m.match_number_in_round, m.round,
			   m.team1_id, t1.name, m.team2_id, t2.name,
			   m.team1_score, m.team2_score, m.winner_team_id, m.status, m.next_match_id, tr.name as tournament_name
		FROM matches m
		JOIN tournaments tr ON m.tournament_id = tr.id
		LEFT JOIN teams t1 ON m.team1_id = t1.id
		LEFT JOIN teams t2 ON m.team2_id = t2.id
		WHERE tr.sport = ?
		ORDER BY m.round, m.match_number_in_round`
	rows, err := r.DB.Query(query, sport)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.MatchResponse
	for rows.Next() {
		var res model.MatchResponse
		var team1Name, team2Name, status, tournamentName sql.NullString
		var team1Score, team2Score sql.NullInt64

		err := rows.Scan(
			&res.ID, &res.MatchNumberInRound, &res.Round,
			&res.Team1ID, &team1Name, &res.Team2ID, &team2Name,
			&team1Score, &team2Score, &res.WinnerTeamID, &status, &res.NextMatchID, &tournamentName,
		)
		if err != nil {
			return nil, err
		}

		// sql.Null* 型からポインタ型へ変換
		if team1Name.Valid {
			res.Team1Name = &team1Name.String
		}
		if team2Name.Valid {
			res.Team2Name = &team2Name.String
		}
		if status.Valid {
			res.Status = &status.String
		}
		if tournamentName.Valid {
			res.TournamentName = &tournamentName.String
		}
		if team1Score.Valid {
			score := int(team1Score.Int64)
			res.Team1Score = &score
		}
		if team2Score.Valid {
			score := int(team2Score.Int64)
			res.Team2Score = &score
		}

		results = append(results, res)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return []model.MatchResponse{}, nil
	}

	return results, nil
}

func (r *MatchRepository) ResetMatch(matchID int, user interface{}) error {
	userMap, ok := user.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unauthorized")
	}
	role, _ := userMap["role"].(string)

	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			log.Printf("Rolling back transaction due to error: %v", err)
			tx.Rollback()
		}
	}()

	var matchSport string
	var previousWinnerID sql.NullInt64
	query := `
		SELECT t.sport, m.winner_team_id
		FROM matches m JOIN tournaments t ON m.tournament_id = t.id
		WHERE m.id = ?`
	err = tx.QueryRow(query, matchID).Scan(&matchSport, &previousWinnerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("not found")
		}
		return fmt.Errorf("failed to query initial match: %w", err)
	}

	if role != "superroot" {
		return fmt.Errorf("forbidden")
	}

	if !previousWinnerID.Valid {
		return fmt.Errorf("match result not submitted yet")
	}

	if err := r.resetMatchResult(tx, matchID); err != nil {
		return fmt.Errorf("failed to reset match result: %w", err)
	}

	// 試合自体の結果をリセット
	_, err = tx.Exec(`UPDATE matches SET team1_score = NULL, team2_score = NULL, winner_team_id = NULL, status = NULL WHERE id = ?`, matchID)
	if err != nil {
		return fmt.Errorf("failed to reset match scores and status: %w", err)
	}

	// 晴天時卓球の場合、雨天時もリセット
	var tournamentName string
	err = tx.QueryRow(`SELECT name FROM tournaments WHERE sport = 'table_tennis' AND id = (SELECT tournament_id FROM matches WHERE id = ?)`, matchID).Scan(&tournamentName)
	if err == nil && tournamentName == "卓球（晴天時）" {
		var rainyMatchID int
		rainyQuery := `
			SELECT m_rainy.id
			FROM matches m_sunny
			JOIN tournaments t_rainy ON t_rainy.name = '卓球（雨天時）'
			JOIN matches m_rainy ON m_rainy.tournament_id = t_rainy.id 
				AND m_rainy.round = m_sunny.round 
				AND m_rainy.match_number_in_round = m_sunny.match_number_in_round
			WHERE m_sunny.id = ?`
		err = tx.QueryRow(rainyQuery, matchID).Scan(&rainyMatchID)
		if err == nil {
			if err := r.resetMatchResult(tx, rainyMatchID); err != nil {
				return fmt.Errorf("failed to reset rainy match result: %w", err)
			}
			_, err = tx.Exec(`UPDATE matches SET team1_score = NULL, team2_score = NULL, winner_team_id = NULL, status = NULL WHERE id = ?`, rainyMatchID)
			if err != nil {
				return fmt.Errorf("failed to reset rainy match scores and status: %w", err)
			}
		}
	}

	return tx.Commit()
}