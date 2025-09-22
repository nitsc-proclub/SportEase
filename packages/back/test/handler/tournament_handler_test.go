package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/saku0512/GYOUJI_HP/backend/internal/handler"
	"github.com/saku0512/GYOUJI_HP/backend/internal/service"
	"github.com/stretchr/testify/assert"
)

type MockTournamentService struct{}

func (m *MockTournamentService) GetTournamentsBySport(sport string) (interface{}, error) {
	return []map[string]interface{}{
		{"tournament": map[string]interface{}{"id": 1, "name": "test", "sport": sport, "weather_condition": "any"}, "teams": []interface{}{}, "matches": []interface{}{}},
	}, nil
}

func TestGetTournamentsBySport(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var svc service.TournamentServiceInterface = &MockTournamentService{}
	h := handler.NewTournamentHandler(svc)
	r := gin.Default()
	r.GET("/api/tournaments/:sport", h.GetTournamentsBySport)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/tournaments/volleyball", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "volleyball")
}
