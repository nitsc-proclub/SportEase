package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/saku0512/GYOUJI_HP/backend/internal/handler"
	"github.com/saku0512/GYOUJI_HP/backend/internal/model"
)

// モックサービス
type mockScoreService struct {
	resp []model.ScoreBreakdown
	err  error
}

func (m *mockScoreService) GetScores() ([]model.ScoreBreakdown, error) { return m.resp, m.err }

func TestScoreHandler_GetScores_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 準備: ハンドラにモックサービスを注入
	svc := &mockScoreService{resp: []model.ScoreBreakdown{
		{
			ClassName:                   "IE4",
			InitScore:                   5,
			AttendanceScore:             8,
			Volleyball1Score:            20,
			VolleyballChampionshipScore: 80,
			TableTennisRainyBonusScore:  10,
			TotalExcludingInit:          118,
			TotalIncludingInit:          123,
		},
	}}
	h := handler.NewScoreHandler(nil)
	// handler.NewScoreHandlerは構造体を返すだけなので、フィールドを直接差し替える
	h.Service = svc

	r := gin.New()
	r.GET("/api/score", h.GetScores)

	req, _ := http.NewRequest(http.MethodGet, "/api/score", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var got []model.ScoreBreakdown
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 item, got %d", len(got))
	}
	if got[0].ClassName != "IE4" || got[0].TotalIncludingInit != 123 {
		t.Fatalf("unexpected payload: %+v", got[0])
	}
}

func TestScoreHandler_GetScores_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// エラーを返すモック
	svc := &mockScoreService{err: assertError("boom")}
	h := handler.NewScoreHandler(nil)
	h.Service = svc

	r := gin.New()
	r.GET("/api/score", h.GetScores)

	req, _ := http.NewRequest(http.MethodGet, "/api/score", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// 簡易エラータイプ
type assertError string

func (e assertError) Error() string { return string(e) }
