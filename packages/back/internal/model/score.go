package model

type ScoreBreakdown struct {
	ClassName                    string `json:"class_name"`
	InitScore                    int    `json:"init_score"`
	AttendanceScore              int    `json:"attendance_score"`
	Volleyball1Score             int    `json:"volleyball1_score"`
	Volleyball2Score             int    `json:"volleyball2_score"`
	Volleyball3Score             int    `json:"volleyball3_score"`
	VolleyballChampionshipScore  int    `json:"volleyball_championship_score"`
	TableTennis1Score            int    `json:"table_tennis1_score"`
	TableTennis2Score            int    `json:"table_tennis2_score"`
	TableTennis3Score            int    `json:"table_tennis3_score"`
	TableTennisChampionshipScore int    `json:"table_tennis_championship_score"`
	TableTennisRainyBonusScore   int    `json:"table_tennis_rainy_bonus_score"`
	Soccer1Score                 int    `json:"soccer1_score"`
	Soccer2Score                 int    `json:"soccer2_score"`
	Soccer3Score                 int    `json:"soccer3_score"`
	SoccerChampionshipScore      int    `json:"soccer_championship_score"`
	RelayAScore                  int    `json:"relay_A_score"`
	RelayBScore                  int    `json:"relay_B_score"`
	RelayBonusScore              int    `json:"relay_bonus_score"`
	TotalExcludingInit           int    `json:"total_excluding_init"`
	TotalIncludingInit           int    `json:"total_including_init"`
}
