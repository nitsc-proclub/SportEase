package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/saku0512/GYOUJI_HP/backend/internal/model"
	"github.com/saku0512/GYOUJI_HP/backend/internal/service"
)

type MatchHandler struct {
	Service service.MatchServiceInterface
}

// GET /api/matches/:sport
// GetMatchesBySport godoc
// @Summary      指定競技の試合一覧取得
// @Description  指定された競技の試合情報を取得します。admin/superrootのみアクセス可能。
// @Tags         matches
// @Produce      json
// @Param        sport path string true "競技名"
// @Param        Authorization header string true "Bearerトークン"
// @Success      200 {array} model.MatchResponse "試合情報一覧"
// @Failure      401 {object} model.ErrorResponse "認証が必要です (Unauthorized)"
// @Failure      403 {object} model.ErrorResponse "権限がありません (Forbidden)"
// @Failure      404 {object} model.ErrorResponse "試合が見つかりません (Not found)"
// @Router       /api/matches/{sport} [get]
// @Security     ApiKeyAuth
func (h *MatchHandler) GetMatchesBySport(c *gin.Context) {
	sport := c.Param("sport")
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userMap, ok := user.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	role, _ := userMap["role"].(string)
	assignedSport, _ := userMap["assigned_sport"].(string)
	// 権限判定
	if !(role == "superroot" || (role == "admin" && assignedSport == sport)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to view matches for this sport."})
		return
	}
	// サービス層から試合一覧取得
	matches, err := h.Service.GetMatchesBySport(sport)
	if err != nil {
		if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "No matches found for this sport."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, matches)
}

func NewMatchHandler(s service.MatchServiceInterface) *MatchHandler {
	return &MatchHandler{Service: s}
}

// PUT /api/matches/:id
// UpdateMatchScore godoc
// @Summary      試合のスコアを更新
// @Description  指定されたIDの試合のスコアを更新し、勝者を次の試合へ進出させます。superrootまたは担当競技のadmin権限が必要です。
// @Tags         matches
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "試合ID"
// @Param        score body      model.TeamScore true  "更新するスコア情報"
// @Param        Authorization header string true "Bearerトークン"
// @Success      200  {object}  model.Match "更新後の試合情報"
// @Failure      400  {object}  model.ErrorResponse "無効なリクエストです (Invalid match ID or request body)"
// @Failure      401  {object}  model.ErrorResponse "認証が必要です (Unauthorized)"
// @Failure      403  {object}  model.ErrorResponse "権限がありません (Forbidden)"
// @Failure      404  {object}  model.ErrorResponse "試合が見つかりません (Match not found)"
// @Failure      500  {object}  model.ErrorResponse "サーバー内部エラー"
// @Router       /api/matches/{id} [put]
// @Security     ApiKeyAuth
func (h *MatchHandler) UpdateMatchScore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	var req model.TeamScore
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, exists := c.Get("user") // AuthMiddlewareでセットされている前提
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	result, err := h.Service.UpdateMatchScore(id, *req.Team1Score, *req.Team2Score, user)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this match."})
			return
		} else if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Match not found."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// POST /api/matches/:id/reset
// ResetMatch godoc
// @Summary      試合結果をリセット
// @Description  指定されたIDの試合結果をリセットします。superrootのみアクセス可能。
// @Tags         matches
// @Produce      json
// @Param        id   path      int  true  "試合ID"
// @Param        Authorization header string true "Bearerトークン"
// @Success      200  {object}  model.SuccessResponse "試合結果が正常にリセットされました。"
// @Failure      400  {object}  model.ErrorResponse "無効なリクエストです (Invalid match ID)"
// @Failure      401  {object}  model.ErrorResponse "認証が必要です (Unauthorized)"
// @Failure      403  {object}  model.ErrorResponse "権限がありません (Forbidden)"
// @Failure      404  {object}  model.ErrorResponse "試合が見つかりません (Match not found)"
// @Failure      500  {object}  model.ErrorResponse "サーバー内部エラー"
// @Router       /api/matches/{id}/reset [post]
// @Security     ApiKeyAuth
func (h *MatchHandler) ResetMatch(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.Service.ResetMatch(id, user); err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to reset this match."})
			return
		} else if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Match not found."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Match result reset successfully"})
}
