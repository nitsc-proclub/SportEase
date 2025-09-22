package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/saku0512/GYOUJI_HP/backend/internal/handler"
	"github.com/saku0512/GYOUJI_HP/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock Service ---

// MockMatchService は service.MatchServiceInterface のモック実装です。
type MockMatchService struct {
	mock.Mock
}

// UpdateMatchScore はインターフェースの要件を満たすモックメソッドです。
func (m *MockMatchService) UpdateMatchScore(matchID int, team1Score int, team2Score int, user interface{}) (interface{}, error) {
	args := m.Called(matchID, team1Score, team2Score, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

// GetMatchesBySport はインターフェースの要件を満たすモックメソッドです。
func (m *MockMatchService) GetMatchesBySport(sport string) ([]model.MatchResponse, error) {
	args := m.Called(sport)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	// 型アサーションを安全に行う
	if val, ok := args.Get(0).([]model.MatchResponse); ok {
		return val, args.Error(1)
	}
	// モックの設定ミスなどで型が違う場合にnilを返す
	return nil, args.Error(1)
}

// --- Handler Tests ---

func TestUpdateMatchScoreHandler(t *testing.T) {
	// Ginをテストモードに設定
	gin.SetMode(gin.TestMode)

	// テスト用のユーザー情報
	testUser := map[string]interface{}{"role": "admin", "assigned_sport": "volleyball"}

	// 期待される成功時のレスポンスデータ
	successResponse := gin.H{
		"ID":           1,
		"Team1Score":   2,
		"Team2Score":   0,
		"WinnerTeamID": 10,
	}

	testCases := []struct {
		name                 string
		matchID              string // URLパラメータは文字列
		requestBody          gin.H  // リクエストボディ
		userContext          interface{}
		setupMock            func(mockService *MockMatchService) // モックの設定
		expectedStatusCode   int
		expectedResponseBody string // 期待されるレスポンスボディ（JSON文字列）
	}{
		{
			name:        "成功: 200 OK",
			matchID:     "1",
			requestBody: gin.H{"team1_score": 2, "team2_score": 0},
			userContext: testUser,
			setupMock: func(mockService *MockMatchService) {
				mockService.On("UpdateMatchScore", 1, 2, 0, testUser).Return(successResponse, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"ID":1,"Team1Score":2,"Team2Score":0,"WinnerTeamID":10}`,
		},
		{
			name:                 "失敗: 400 Bad Request (無効なMatch ID)",
			matchID:              "abc",
			requestBody:          gin.H{"team1_score": 2, "team2_score": 0},
			userContext:          testUser,
			setupMock:            func(mockService *MockMatchService) {}, // サービスは呼ばれない
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Invalid match ID"}`,
		},
		{
			name:                 "失敗: 400 Bad Request (不正なリクエストボディ)",
			matchID:              "1",
			requestBody:          gin.H{"invalid_key": "value"}, // team1_score, team2_score がない
			userContext:          testUser,
			setupMock:            func(mockService *MockMatchService) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Invalid request body"}`,
		},
		{
			name:                 "失敗: 401 Unauthorized (ユーザー情報がコンテキストにない)",
			matchID:              "1",
			requestBody:          gin.H{"team1_score": 2, "team2_score": 0},
			userContext:          nil, // ユーザー情報をnilにする
			setupMock:            func(mockService *MockMatchService) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":"Unauthorized"}`,
		},
		{
			name:        "失敗: 403 Forbidden (サービスから権限エラーが返る)",
			matchID:     "1",
			requestBody: gin.H{"team1_score": 2, "team2_score": 0},
			userContext: testUser,
			setupMock: func(mockService *MockMatchService) {
				mockService.On("UpdateMatchScore", 1, 2, 0, testUser).Return(nil, fmt.Errorf("forbidden"))
			},
			expectedStatusCode:   http.StatusForbidden,
			expectedResponseBody: `{"error":"You do not have permission to update this match."}`,
		},
		{
			name:        "失敗: 404 Not Found (サービスから見つからないエラーが返る)",
			matchID:     "999",
			requestBody: gin.H{"team1_score": 2, "team2_score": 0},
			userContext: testUser,
			setupMock: func(mockService *MockMatchService) {
				mockService.On("UpdateMatchScore", 999, 2, 0, testUser).Return(nil, fmt.Errorf("not found"))
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"error":"Match not found."}`,
		},
		{
			name:        "失敗: 500 Internal Server Error (サービスから予期せぬエラーが返る)",
			matchID:     "1",
			requestBody: gin.H{"team1_score": 2, "team2_score": 0},
			userContext: testUser,
			setupMock: func(mockService *MockMatchService) {
				mockService.On("UpdateMatchScore", 1, 2, 0, testUser).Return(nil, fmt.Errorf("some database error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"some database error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// モックサービスとハンドラのインスタンスを作成
			mockService := new(MockMatchService)
			tc.setupMock(mockService)
			handler := handler.NewMatchHandler(mockService)

			// ルーターのセットアップ
			router := gin.Default()
			router.PUT("/api/matches/:id", func(c *gin.Context) {
				// Middlewareの動作をシミュレート
				if tc.userContext != nil {
					c.Set("user", tc.userContext)
				}
				handler.UpdateMatchScore(c)
			})

			// リクエストボディをJSONに変換
			jsonBody, _ := json.Marshal(tc.requestBody)
			bodyReader := bytes.NewReader(jsonBody)

			// HTTPリクエストを作成
			req, _ := http.NewRequest(http.MethodPut, "/api/matches/"+tc.matchID, bodyReader)
			req.Header.Set("Content-Type", "application/json")

			// レスポンスを記録するレコーダーを作成
			w := httptest.NewRecorder()

			// リクエストを実行
			router.ServeHTTP(w, req)

			// 結果を検証
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.JSONEq(t, tc.expectedResponseBody, w.Body.String())

			// モックの期待が満たされたか検証
			mockService.AssertExpectations(t)
		})
	}
}

// TestMatchHandler_GetMatchesBySport は GetMatchesBySport ハンドラをテストします。
func TestMatchHandler_GetMatchesBySport(t *testing.T) {
	gin.SetMode(gin.TestMode)

	teamAName, teamBName := "Team A", "Team B"
	teamCName, teamDName := "Team C", "Team D"
	mockMatches := []model.MatchResponse{
		{ID: 1, Round: 1, Team1Name: &teamAName, Team2Name: &teamBName},
		{ID: 2, Round: 1, Team1Name: &teamCName, Team2Name: &teamDName},
	}
	// 期待されるJSON文字列もこのデータから生成
	mockMatchesJSON, _ := json.Marshal(mockMatches)

	testCases := []struct {
		name                 string
		sportParam           string
		userContext          interface{}
		setupMock            func(mockService *MockMatchService)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:       "正常系: adminユーザーが担当競技の試合を取得",
			sportParam: "volleyball",
			userContext: map[string]interface{}{
				"role":           "admin",
				"assigned_sport": "volleyball",
			},
			setupMock: func(mockService *MockMatchService) {
				mockService.On("GetMatchesBySport", "volleyball").Return(mockMatches, nil).Once()
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: string(mockMatchesJSON), // 期待値も正しいJSONに
		},
		{
			name:       "正常系: superrootユーザーが試合を取得",
			sportParam: "soccer",
			userContext: map[string]interface{}{
				"role":           "superroot",
				"assigned_sport": "",
			},
			setupMock: func(mockService *MockMatchService) {
				mockService.On("GetMatchesBySport", "soccer").Return(mockMatches, nil).Once()
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: string(mockMatchesJSON), // 期待値も正しいJSONに
		},
		{
			name:       "異常系: adminユーザーが担当外の競技で403 Forbidden",
			sportParam: "soccer",
			userContext: map[string]interface{}{
				"role":           "admin",
				"assigned_sport": "volleyball",
			},
			setupMock:            func(mockService *MockMatchService) {}, // サービスは呼ばれない
			expectedStatusCode:   http.StatusForbidden,
			expectedResponseBody: `{"error":"You do not have permission to view matches for this sport."}`,
		},
		{
			name:                 "異常系: ユーザー情報がなく401 Unauthorized",
			sportParam:           "volleyball",
			userContext:          nil,
			setupMock:            func(mockService *MockMatchService) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":"Unauthorized"}`,
		},
		{
			name:       "異常系: サービスからnot foundエラーが返され404 Not Found",
			sportParam: "tennis",
			userContext: map[string]interface{}{
				"role": "superroot",
			},
			setupMock: func(mockService *MockMatchService) {
				mockService.On("GetMatchesBySport", "tennis").Return(([]model.MatchResponse)(nil), fmt.Errorf("not found")).Once()
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"error":"No matches found for this sport."}`,
		},
		{
			name:       "異常系: サービスから予期せぬエラーが返され500 Internal Server Error",
			sportParam: "volleyball",
			userContext: map[string]interface{}{
				"role": "superroot",
			},
			setupMock: func(mockService *MockMatchService) {
				mockService.On("GetMatchesBySport", "volleyball").Return(([]model.MatchResponse)(nil), fmt.Errorf("database connection failed")).Once()
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"database connection failed"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockMatchService)
			tc.setupMock(mockService)
			handler := handler.NewMatchHandler(mockService)

			router := gin.Default()
			router.GET("/api/matches/:sport", func(c *gin.Context) {
				if tc.userContext != nil {
					c.Set("user", tc.userContext)
				}
				handler.GetMatchesBySport(c)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/matches/"+tc.sportParam, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.JSONEq(t, tc.expectedResponseBody, w.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}
