package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saku0512/GYOUJI_HP/backend/internal/model"
	"github.com/saku0512/GYOUJI_HP/backend/internal/service"
)

type RelayHandler struct {
	service service.RelayService
}

func NewRelayHandler(s service.RelayService) *RelayHandler {
	return &RelayHandler{service: s}
}

// GetRelayRankings リレーの指定されたブロックの順位を取得
// @Summary リレーブロックの順位取得
// @Description 指定されたブロック（A or B）のリレー順位を取得します
// @Tags relay
// @Param block query string true "リレーブロック (A or B)"
// @Success 200 {object} model.RelayRankResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/relay [get]
func (h *RelayHandler) GetRelayRankings(c *gin.Context) {
	block := c.Query("block")
	if block == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "block parameter is required"})
		return
	}

	rankings, err := h.service.GetRelayRankings(block)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := model.RelayRankResponse{
		Block:    block,
		Rankings: rankings,
	}

	c.JSON(http.StatusOK, response)
}

// RegisterRelayRankings リレーの指定されたブロックの順位を登録
// @Summary リレーブロックの順位登録
// @Description 指定されたブロック（A or B）のリレー順位を登録します
// @Tags relay
// @Param block query string true "リレーブロック (A or B)"
// @Param body body model.RelayRankRequest true "順位データ"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/relay [post]
func (h *RelayHandler) RegisterRelayRankings(c *gin.Context) {
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
	username, _ := userMap["username"].(string)
	assignedSport, _ := userMap["assigned_sport"].(string)

	// 権限チェック（superrootまたはadmin_relayのみ）
	if !(role == "superroot" || (role == "admin" && username == "admin_relay" && assignedSport == "relay")) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: only superroot or admin_relay can register relay rankings"})
		return
	}

	// ブロックパラメータの取得
	block := c.Query("block")
	if block == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "block parameter is required"})
		return
	}

	// リクエストボディの解析
	var req model.RelayRankRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// 順位を登録
	err := h.service.RegisterRelayRankings(block, req.Rankings)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "relay rankings registered successfully"})
}