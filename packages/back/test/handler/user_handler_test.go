package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/saku0512/GYOUJI_HP/backend/internal/handler"
	"github.com/saku0512/GYOUJI_HP/backend/internal/router"
	"github.com/saku0512/GYOUJI_HP/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- モックの定義 ---

// MockUserService は service.UserService のインターフェースを満たすモックオブジェクトです。
// これにより、実際のビジネスロジック（DBアクセスやJWT生成）をシミュレートします。
type MockUserService struct {
	mock.Mock
}

// Login はモックされたログインメソッドです。
// テストケースごとに「この引数で呼ばれたら、この値を返す」という振る舞いを定義できます。
func (m *MockUserService) Login(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

// --- テスト本体 ---

func TestUserHandler_Login(t *testing.T) {
	// Ginをテストモードに設定
	gin.SetMode(gin.TestMode)

	// テストケースを定義
	testCases := []struct {
		name            string            // テストケース名
		requestBody     map[string]string // リクエストボディ
		mockUsername    string            // モックに渡すユーザー名
		mockPassword    string            // モックに渡すパスワード
		mockReturnToken string            // モックが返すトークン
		mockReturnError error             // モックが返すエラー
		expectedStatus  int               // 期待するHTTPステータスコード
		expectedBody    string            // 期待するレスポンスボディ
	}{
		{
			name:            "正常系: ログイン成功",
			requestBody:     map[string]string{"username": "testuser", "password": "password123"},
			mockUsername:    "testuser",
			mockPassword:    "password123",
			mockReturnToken: "mocked.jwt.token",
			mockReturnError: nil,
			expectedStatus:  http.StatusOK,
			expectedBody:    `{"token":"mocked.jwt.token"}`,
		},
		{
			name:            "異常系: 認証情報が無効",
			requestBody:     map[string]string{"username": "wronguser", "password": "wrongpassword"},
			mockUsername:    "wronguser",
			mockPassword:    "wrongpassword",
			mockReturnToken: "",
			mockReturnError: errors.New("invalid credentials"),
			expectedStatus:  http.StatusUnauthorized,
			expectedBody:    `{"error":"Invalid username or password"}`,
		},
		{
			name:           "異常系: リクエストボディが不正 (パスワード欠け)",
			requestBody:    map[string]string{"username": "testuser"},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid request body"}`,
		},
		{
			name:           "異常系: リクエストボディが空",
			requestBody:    map[string]string{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid request body"}`,
		},
	}

	// 各テストケースを実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// --- セットアップ ---
			mockService := new(MockUserService)

			// モックの振る舞いを設定（リクエストボディが正常な場合のみ）
			if tc.expectedStatus != http.StatusBadRequest {
				mockService.On("Login", tc.mockUsername, tc.mockPassword).Return(tc.mockReturnToken, tc.mockReturnError)
			}

			// ハンドラとルーターを初期化
			userHandler := handler.NewUserHandler(mockService)
			router := gin.Default()
			router.POST("/login", userHandler.Login)

			// リクエストボディをJSONに変換
			body, _ := json.Marshal(tc.requestBody)

			// HTTPリクエストをシミュレート
			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// レスポンスを記録するためのレコーダーを作成
			w := httptest.NewRecorder()

			// --- 実行 ---
			router.ServeHTTP(w, req)

			// --- 検証 ---
			// ステータスコードを検証
			assert.Equal(t, tc.expectedStatus, w.Code)
			// レスポンスボディを検証
			assert.JSONEq(t, tc.expectedBody, w.Body.String())

			// モックが期待通りに呼ばれたか検証
			if tc.expectedStatus != http.StatusBadRequest {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestUserHandler_GetMe(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// --- テストケースの定義 ---
	testCases := []struct {
		name                 string      // テストケース名
		userContext          interface{} // Ginコンテキストにセットするユーザー情報
		expectedStatusCode   int         // 期待されるHTTPステータスコード
		expectedResponseBody string      // 期待されるレスポンスボディ(JSON文字列)
	}{
		{
			name: "正常系: ユーザー情報がコンテキストに存在する",
			userContext: map[string]interface{}{
				"user_id":        1,
				"username":       "admin",
				"role":           "admin",
				"assigned_sport": "soccer",
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"user_id":1, "username":"admin", "role":"admin", "assigned_sport":"soccer"}`,
		},
		{
			name:                 "異常系: ユーザー情報がコンテキストに存在しない",
			userContext:          nil, // ユーザー情報がないケース
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":"Unauthorized"}`,
		},
	}

	// --- テストの実行 ---
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// モックサービスとハンドラの準備
			mockService := new(MockUserService) // GetMeでは使われないが初期化に必要
			userHandler := handler.NewUserHandler(mockService)

			// レスポンスレコーダーとルーターの準備
			w := httptest.NewRecorder()
			router := gin.Default()

			// テスト対象のエンドポイントを定義
			router.GET("/api/auth/me", func(c *gin.Context) {
				// ミドルウェアの動作をシミュレート
				if tc.userContext != nil {
					c.Set("user", tc.userContext)
				}
				userHandler.GetMe(c)
			})

			// HTTPリクエストの作成と実行
			req, _ := http.NewRequest(http.MethodGet, "/api/auth/me", nil)
			router.ServeHTTP(w, req)

			// 結果の検証
			assert.Equal(t, tc.expectedStatusCode, w.Code, "ステータスコードが期待値と一致しません")

			// assert.JSONEq を使うことで、キーの順序が違っても内容が同じであればパスします。
			assert.JSONEq(t, tc.expectedResponseBody, w.Body.String(), "レスポンスボディが期待値と一致しません")
		})
	}
}

func TestUserHandler_Logout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const jwtSecret = "test_secret_key"

	// テスト用の有効なJWTを生成
	claims := &service.Claims{
		UserID:   1,
		Username: "testuser",
		Role:     "student",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	validToken, _ := token.SignedString([]byte(jwtSecret))

	testCases := []struct {
		name           string
		authHeader     string // Authorizationヘッダーの値
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "正常系: ログアウト成功",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Logout successful"}`,
		},
		{
			name:           "異常系: Authorizationヘッダーなし",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Authorization header is required"}`,
		},
		{
			name:           "異常系: Bearerプレフィックスなし",
			authHeader:     validToken,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Authorization header format must be Bearer {token}"}`,
		},
		{
			name:           "異常系: 不正なトークン",
			authHeader:     "Bearer invalid.token",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Invalid or expired token"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// --- セットアップ ---
			// ログアウトハンドラはServiceを呼ばないので、モックは不要
			mockService := new(MockUserService)
			userHandler := handler.NewUserHandler(mockService)

			// ルーターをセットアップ（ミドルウェアを有効にするためjwtSecretを渡す）
			r := router.SetupRouter(userHandler, nil, nil, jwtSecret, nil, nil, nil, nil)

			// HTTPリクエストをシミュレート
			req, _ := http.NewRequest(http.MethodPost, "/api/auth/logout", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}

			w := httptest.NewRecorder()

			// --- 実行 ---
			r.ServeHTTP(w, req)

			// --- 検証 ---
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}
