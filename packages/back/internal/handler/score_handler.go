package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saku0512/GYOUJI_HP/backend/internal/service"
)

type ScoreHandler struct{ Service service.ScoreServiceInterface }

func NewScoreHandler(s service.ScoreServiceInterface) *ScoreHandler { return &ScoreHandler{Service: s} }

// GetScores 全クラスの出席点を取得
// @Summary 出席点一覧取得
// @Description 全クラスの出席点を取得します
// @Tags score
// @Produce json
// @Success 200 {array} model.ScoreBreakdown
// @Failure 500 {object} map[string]string
// @Router /api/score [get]
// @Security ApiKeyAuth
func (h *ScoreHandler) GetScores(c *gin.Context) {
	scores, err := h.Service.GetScores()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, scores)
}
