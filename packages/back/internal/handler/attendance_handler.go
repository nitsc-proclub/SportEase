package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/saku0512/GYOUJI_HP/backend/internal/model"
	"github.com/saku0512/GYOUJI_HP/backend/internal/service"
)

type AttendanceHandler struct {
	service service.AttendanceService
}

func NewAttendanceHandler(s service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{service: s}
}

// GetAttendanceScores 全クラスの出席点を取得
// @Summary 出席点一覧取得
// @Description 全クラスの出席点を取得します
// @Tags attendance
// @Produce json
// @Success 200 {array} model.AttendanceScore
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/attendance [get]
// @Security ApiKeyAuth
func (h *AttendanceHandler) GetAttendanceScores(c *gin.Context) {
	// 認証チェック
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userMap, ok := user.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role, _ := userMap["role"].(string)

	// superrootのみアクセス可能
	if role != "superroot" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: only superroot can access attendance scores"})
		return
	}

	scores, err := h.service.GetAttendanceScores()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, scores)
}

// UpdateAttendanceScore 出席点を更新
// @Summary 出席点更新
// @Description 指定されたクラスの出席点を更新します
// @Tags attendance
// @Accept json
// @Produce json
// @Param class_id path int true "クラスID"
// @Param body body model.AttendanceUpdateRequest true "出席点データ"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/attendance/{class_id} [put]
// @Security ApiKeyAuth
func (h *AttendanceHandler) UpdateAttendanceScore(c *gin.Context) {
	// 認証チェック
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userMap, ok := user.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role, _ := userMap["role"].(string)

	// superrootのみアクセス可能
	if role != "superroot" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: only superroot can update attendance scores"})
		return
	}

	// クラスIDの取得
	classIDStr := c.Param("class_id")
	classID, err := strconv.Atoi(classIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid class_id"})
		return
	}

	// リクエストボディの解析
	var req struct {
		Score int `json:"score" binding:"required,min=0,max=10"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body or score out of range (0-10)"})
		return
	}

	// 出席点を更新
	err = h.service.UpdateAttendanceScore(classID, req.Score)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "attendance score updated successfully"})
}

// BatchUpdateAttendanceScores 出席点を一括更新
// @Summary 出席点一括更新
// @Description 複数クラスの出席点を一括で更新します
// @Tags attendance
// @Accept json
// @Produce json
// @Param body body model.BatchAttendanceUpdateRequest true "出席点データ"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/attendance/batch [put]
// @Security ApiKeyAuth
func (h *AttendanceHandler) BatchUpdateAttendanceScores(c *gin.Context) {
	// 認証チェック
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userMap, ok := user.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role, _ := userMap["role"].(string)

	// superrootのみアクセス可能
	if role != "superroot" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: only superroot can update attendance scores"})
		return
	}

	// リクエストボディの解析
	var req model.BatchAttendanceUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// 出席点を一括更新
	err := h.service.BatchUpdateAttendanceScores(req.Scores)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "attendance scores updated successfully"})
}