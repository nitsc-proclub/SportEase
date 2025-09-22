package repository

import (
	"database/sql"
	"fmt"

	"github.com/saku0512/GYOUJI_HP/backend/internal/model"
)

type RelayRepository struct {
	DB *sql.DB
}

func NewRelayRepository(db *sql.DB) *RelayRepository {
	return &RelayRepository{DB: db}
}

// GetRelayRankings 指定されたブロックのリレー順位を取得
func (r *RelayRepository) GetRelayRankings(block string) (map[int]int, error) {
	query := `
		SELECT relay_rank, class_id 
		FROM relay_results 
		WHERE relay_type = ? AND class_id > 0
		ORDER BY relay_rank ASC
	`
	
	rows, err := r.DB.Query(query, block)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rankings := make(map[int]int) // 順位 -> 学年
	
	for rows.Next() {
		var rank, classID int
		if err := rows.Scan(&rank, &classID); err != nil {
			return nil, err
		}
		
		// クラスIDから学年を逆算
		grade := r.classIDToGrade(classID)
		if grade > 0 {
			rankings[rank] = grade
		}
	}
	
	return rankings, nil
}

// RegisterRelayRankings 指定されたブロックのリレー順位を登録
func (r *RelayRepository) RegisterRelayRankings(block string, rankings map[int]int) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 該当ブロックのリレー得点をリセット
	scoreColumn := fmt.Sprintf("relay_%s_score", block)
	_, err = tx.Exec(fmt.Sprintf("UPDATE team_points SET %s = 0", scoreColumn))
	if err != nil {
		return err
	}

	// 順位ごとに結果を登録
	for rank, grade := range rankings {
		// relay_resultsテーブルを更新（代表クラスIDを使用）
		representativeClassID := r.getRepresentativeClassID(grade)
		_, err = tx.Exec(
			"UPDATE relay_results SET class_id = ? WHERE relay_type = ? AND relay_rank = ?",
			representativeClassID, block, rank,
		)
		if err != nil {
			return err
		}

		// 該当学年の全クラスに得点を付与
		score := model.RelayScoreConfig[rank]
		classIDs := model.GradeToClassIDs[grade]
		
		for _, classID := range classIDs {
			_, err = tx.Exec(
				fmt.Sprintf("UPDATE team_points SET %s = ? WHERE class_id = ?", scoreColumn),
				score, classID,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

// GetGradeScores 全学年のリレー得点を取得
func (r *RelayRepository) GetGradeScores() ([]model.GradeScore, error) {
	query := `
		SELECT 
			CASE 
				WHEN class_id BETWEEN 11 AND 13 THEN 1
				WHEN class_id BETWEEN 21 AND 23 THEN 2
				WHEN class_id BETWEEN 31 AND 33 THEN 3
				WHEN class_id BETWEEN 41 AND 43 THEN 4
				WHEN class_id BETWEEN 51 AND 53 THEN 5
				WHEN class_id = 6 THEN 6
			END as grade,
			MAX(relay_A_score) as block_a_score,
			MAX(relay_B_score) as block_b_score,
			MAX(relay_bonus_score) as bonus_score
		FROM team_points 
		WHERE class_id IN (11, 12, 13, 21, 22, 23, 31, 32, 33, 41, 42, 43, 51, 52, 53, 6)
		GROUP BY grade
		ORDER BY grade
	`
	
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gradeScores []model.GradeScore
	for rows.Next() {
		var gs model.GradeScore
		if err := rows.Scan(&gs.Grade, &gs.BlockAScore, &gs.BlockBScore, &gs.BonusScore); err != nil {
			return nil, err
		}
		gs.TotalScore = gs.BlockAScore + gs.BlockBScore
		gradeScores = append(gradeScores, gs)
	}
	
	return gradeScores, nil
}

// UpdateRelayBonusScores リレー最終順位ボーナス得点を更新
func (r *RelayRepository) UpdateRelayBonusScores(gradeScores []model.GradeScore) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// ボーナス得点をリセット
	_, err = tx.Exec("UPDATE team_points SET relay_bonus_score = 0")
	if err != nil {
		return err
	}

	// 各学年のボーナス得点を更新
	for _, gs := range gradeScores {
		if gs.BonusScore > 0 {
			classIDs := model.GradeToClassIDs[gs.Grade]
			for _, classID := range classIDs {
				_, err = tx.Exec(
					"UPDATE team_points SET relay_bonus_score = ? WHERE class_id = ?",
					gs.BonusScore, classID,
				)
				if err != nil {
					return err
				}
			}
		}
	}

	return tx.Commit()
}

// CheckBothBlocksCompleted 両ブロックの結果が揃っているかチェック
func (r *RelayRepository) CheckBothBlocksCompleted() (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM relay_results 
		WHERE class_id > 0
	`
	
	var count int
	err := r.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return false, err
	}
	
	// 両ブロック（A, B）× 6順位 = 12レコードが必要
	return count >= 12, nil
}

// classIDToGrade クラスIDから学年を取得
func (r *RelayRepository) classIDToGrade(classID int) int {
	for grade, classIDs := range model.GradeToClassIDs {
		for _, id := range classIDs {
			if id == classID {
				return grade
			}
		}
	}
	return 0
}

// getRepresentativeClassID 学年の代表クラスIDを取得
func (r *RelayRepository) getRepresentativeClassID(grade int) int {
	classIDs := model.GradeToClassIDs[grade]
	if len(classIDs) > 0 {
		return classIDs[0] // 最初のクラスIDを代表として使用
	}
	return 0
}