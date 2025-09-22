package model

// User はデータベースの users テーブルに対応する構造体です。
type User struct {
	ID             int64  `json:"id"`
	Username       string `json:"username"`
	HashedPassword string `json:"-"` // パスワードはJSONに含めない
	Role           string `json:"role"`
	AssignedSport  string `json:"assigned_sport"`
}
