package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saku0512/GYOUJI_HP/backend/internal/service"
)

type TournamentHandler struct {
	Service service.TournamentServiceInterface
}

func NewTournamentHandler(s service.TournamentServiceInterface) *TournamentHandler {
	return &TournamentHandler{Service: s}
}

// GET /api/tournaments/:sport
// GetTournamentsBySport godoc
// @Summary      指定競技のトーナメント情報を取得
// @Description  競技名をパスパラメータとして受け取り、関連するトーナメント（晴天時・雨天時など）の情報を返します。
// @Tags         Tournaments
// @Accept       json
// @Produce      json
// @Param        sport   path      string  true  "競技名 (volleyball, table_tennis, soccer)"
// @Success      200  {object}  interface{}
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/tournaments/{sport} [get]
func (h *TournamentHandler) GetTournamentsBySport(c *gin.Context) {
	sport := c.Param("sport")
	result, err := h.Service.GetTournamentsBySport(sport)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
