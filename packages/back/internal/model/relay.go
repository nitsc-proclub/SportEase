package model

import "time"

// RelayBlock リレーブロック情報
type RelayBlock struct {
	Block string `json:"block"` // "A" or "B"
	Name  string `json:"name"`  // "Aブロック" or "Bブロック"
}

// RelayRankRequest リレー順位登録リクエスト
type RelayRankRequest struct {
	Rankings map[int]int `json:"rankings"` // 順位 -> 学年 (1->3, 2->1, ...)
}

// RelayRankResponse リレー順位レスポンス
type RelayRankResponse struct {
	Block    string      `json:"block"`
	Rankings map[int]int `json:"rankings"` // 順位 -> 学年
}

// RelayResult リレー結果（内部処理用）
type RelayResult struct {
	ID        int       `json:"id"`
	Block     string    `json:"block"`     // "A" or "B"
	Rank      int       `json:"rank"`      // 1-6
	Grade     int       `json:"grade"`     // 学年 (1-5, 6=専・教)
	Score     int       `json:"score"`     // 獲得得点
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GradeScore 学年別得点集計
type GradeScore struct {
	Grade      int `json:"grade"`       // 学年
	BlockAScore int `json:"block_a_score"` // Aブロック得点
	BlockBScore int `json:"block_b_score"` // Bブロック得点
	TotalScore int `json:"total_score"`   // 合計得点
	FinalRank  int `json:"final_rank"`    // 最終順位
	BonusScore int `json:"bonus_score"`   // ボーナス得点
}

// RelayScoreConfig リレー得点設定
var RelayScoreConfig = map[int]int{
	1: 30, // 1位: 30点
	2: 25, // 2位: 25点
	3: 20, // 3位: 20点
	4: 15, // 4位: 15点
	5: 10, // 5位: 10点
	6: 5,  // 6位: 5点
}

// RelayBonusConfig リレー最終順位ボーナス得点設定
var RelayBonusConfig = map[int]int{
	1: 30, // 1位: 30点
	2: 20, // 2位: 20点
	3: 10, // 3位: 10点
}

// GradeToClassIDs 学年からクラスIDへのマッピング
var GradeToClassIDs = map[int][]int{
	1: {11, 12, 13}, // 1年生: 1-1, 1-2, 1-3
	2: {21, 22, 23}, // 2年生: IS2, IT2, IE2
	3: {31, 32, 33}, // 3年生: IS3, IT3, IE3
	4: {41, 42, 43}, // 4年生: IS4, IT4, IE4
	5: {51, 52, 53}, // 5年生: IS5, IT5, IE5
	6: {6},          // 専・教: 専・教
}