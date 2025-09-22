package repository

import (
	"database/sql"

	"github.com/saku0512/GYOUJI_HP/backend/internal/model"
)

type ScoreRepository struct {
	DB *sql.DB
}

func NewScoreRepository(db *sql.DB) *ScoreRepository { return &ScoreRepository{DB: db} }

func (r *ScoreRepository) GetScores() ([]model.ScoreBreakdown, error) {
	// チームごとのポイントをteam_pointsテーブルの各カラムで集計
	query := `
	 SELECT t.name AS class_name,
		 COALESCE(MAX(tp.init_score),0) AS init_score,
		 COALESCE(MAX(tp.attendance_score),0) AS attendance_score,
		 COALESCE(MAX(tp.volleyball1_score),0) AS volleyball1_score,
		 COALESCE(MAX(tp.volleyball2_score),0) AS volleyball2_score,
		 COALESCE(MAX(tp.volleyball3_score),0) AS volleyball3_score,
		 COALESCE(MAX(tp.volleyball_championship_score),0) AS volleyball_championship_score,
		 COALESCE(MAX(tp.table_tennis1_score),0) AS table_tennis1_score,
		 COALESCE(MAX(tp.table_tennis2_score),0) AS table_tennis2_score,
		 COALESCE(MAX(tp.table_tennis3_score),0) AS table_tennis3_score,
		 COALESCE(MAX(tp.table_tennis_championship_score),0) AS table_tennis_championship_score,
		 COALESCE(MAX(tp.table_tennis_rainy_bonus_score),0) AS table_tennis_rainy_bonus_score,
		 COALESCE(MAX(tp.soccer1_score),0) AS soccer1_score,
		 COALESCE(MAX(tp.soccer2_score),0) AS soccer2_score,
		 COALESCE(MAX(tp.soccer3_score),0) AS soccer3_score,
		 COALESCE(MAX(tp.soccer_championship_score),0) AS soccer_championship_score,
		 COALESCE(MAX(tp.relay_A_score),0) AS relay_A_score,
		 COALESCE(MAX(tp.relay_B_score),0) AS relay_B_score,
		 COALESCE(MAX(tp.relay_bonus_score),0) AS relay_bonus_score
	 FROM teams t
	 LEFT JOIN team_points tp ON tp.class_id = t.class_id
	 GROUP BY t.name
	 ORDER BY (COALESCE(MAX(tp.init_score),0)
		 + COALESCE(MAX(tp.volleyball1_score),0)
		 + COALESCE(MAX(tp.volleyball2_score),0)
		 + COALESCE(MAX(tp.volleyball3_score),0)
		 + COALESCE(MAX(tp.volleyball_championship_score),0)
		 + COALESCE(MAX(tp.table_tennis1_score),0)
		 + COALESCE(MAX(tp.table_tennis2_score),0)
		 + COALESCE(MAX(tp.table_tennis3_score),0)
		 + COALESCE(MAX(tp.table_tennis_championship_score),0)
		 + COALESCE(MAX(tp.table_tennis_rainy_bonus_score),0)
		 + COALESCE(MAX(tp.soccer1_score),0)
		 + COALESCE(MAX(tp.soccer2_score),0)
		 + COALESCE(MAX(tp.soccer3_score),0)
		 + COALESCE(MAX(tp.soccer_championship_score),0)
		 + COALESCE(MAX(tp.relay_A_score),0)
		 + COALESCE(MAX(tp.relay_B_score),0)
		 + COALESCE(MAX(tp.relay_bonus_score),0)) DESC, t.name ASC
	    `
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.ScoreBreakdown
	for rows.Next() {
		var s model.ScoreBreakdown
		if err := rows.Scan(
			&s.ClassName,
			&s.InitScore,
			&s.AttendanceScore,
			&s.Volleyball1Score,
			&s.Volleyball2Score,
			&s.Volleyball3Score,
			&s.VolleyballChampionshipScore,
			&s.TableTennis1Score,
			&s.TableTennis2Score,
			&s.TableTennis3Score,
			&s.TableTennisChampionshipScore,
			&s.TableTennisRainyBonusScore,
			&s.Soccer1Score,
			&s.Soccer2Score,
			&s.Soccer3Score,
			&s.SoccerChampionshipScore,
			&s.RelayAScore,
			&s.RelayBScore,
			&s.RelayBonusScore,
		); err != nil {
			return nil, err
		}
		s.TotalExcludingInit = s.AttendanceScore + s.Volleyball1Score + s.Volleyball2Score + s.Volleyball3Score + s.VolleyballChampionshipScore + s.TableTennis1Score + s.TableTennis2Score + s.TableTennis3Score + s.TableTennisChampionshipScore + s.TableTennisRainyBonusScore + s.Soccer1Score + s.Soccer2Score + s.Soccer3Score + s.SoccerChampionshipScore + s.RelayAScore + s.RelayBScore + s.RelayBonusScore
		s.TotalIncludingInit = s.InitScore + s.TotalExcludingInit
		results = append(results, s)
	}
	return results, nil
}
