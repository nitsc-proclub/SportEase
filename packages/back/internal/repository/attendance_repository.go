package repository

import (
	"database/sql"

	"github.com/saku0512/GYOUJI_HP/backend/internal/model"
)

type AttendanceRepository struct {
	DB *sql.DB
}

func NewAttendanceRepository(db *sql.DB) *AttendanceRepository {
	return &AttendanceRepository{DB: db}
}

// GetAttendanceScores 全クラスの出席点を取得
func (r *AttendanceRepository) GetAttendanceScores() ([]model.AttendanceScore, error) {
	query := `
		SELECT 
			tp.class_id,
			t.name as class_name,
			COALESCE(tp.attendance_score, 0) as score
		FROM team_points tp
		JOIN teams t ON t.class_id = tp.class_id
		WHERE tp.class_id IN (11, 12, 13, 21, 22, 23, 31, 32, 33, 41, 42, 43, 51, 52, 53, 6)
		GROUP BY tp.class_id, t.name
		ORDER BY tp.class_id
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []model.AttendanceScore
	for rows.Next() {
		var score model.AttendanceScore
		if err := rows.Scan(&score.ClassID, &score.ClassName, &score.Score); err != nil {
			return nil, err
		}
		scores = append(scores, score)
	}

	return scores, nil
}

// UpdateAttendanceScore 出席点を更新
func (r *AttendanceRepository) UpdateAttendanceScore(classID int, score int) error {
	// team_pointsレコードが存在しない場合は作成
	_, err := r.DB.Exec(`
		INSERT INTO team_points (class_id, attendance_score) 
		VALUES (?, ?) 
		ON DUPLICATE KEY UPDATE attendance_score = VALUES(attendance_score)
	`, classID, score)

	return err
}

// BatchUpdateAttendanceScores 出席点を一括更新
func (r *AttendanceRepository) BatchUpdateAttendanceScores(scores []model.AttendanceUpdateItem) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO team_points (class_id, attendance_score) 
		VALUES (?, ?) 
		ON DUPLICATE KEY UPDATE attendance_score = VALUES(attendance_score)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range scores {
		_, err = stmt.Exec(item.ClassID, item.Score)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}