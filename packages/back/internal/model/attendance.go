package model

// AttendanceScore 出席点情報
type AttendanceScore struct {
	ClassID   int    `json:"class_id"`
	ClassName string `json:"class_name"`
	Score     int    `json:"score"`
}

// AttendanceUpdateRequest 出席点更新リクエスト
type AttendanceUpdateRequest struct {
	Score int `json:"score" binding:"required,min=0,max=10"`
}

// BatchAttendanceUpdateRequest 出席点一括更新リクエスト
type BatchAttendanceUpdateRequest struct {
	Scores []AttendanceUpdateItem `json:"scores" binding:"required"`
}

// AttendanceUpdateItem 出席点更新アイテム
type AttendanceUpdateItem struct {
	ClassID int `json:"class_id" binding:"required"`
	Score   int `json:"score" binding:"required,min=0,max=10"`
}