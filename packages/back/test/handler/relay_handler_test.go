package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/saku0512/GYOUJI_HP/backend/internal/handler"
	"github.com/saku0512/GYOUJI_HP/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRelayService テスト用のモックサービス
type MockRelayService struct {
	mock.Mock
}

func (m *MockRelayService) GetRelayRankings(block string) (map[int]int, error) {
	args := m.Called(block)
	return args.Get(0).(map[int]int), args.Error(1)
}

func (m *MockRelayService) RegisterRelayRankings(block string, rankings map[int]int) error {
	args := m.Called(block, rankings)
	return args.Error(0)
}

func setupRelayRouter(mockService *MockRelayService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	
	relayHandler := handler.NewRelayHandler(mockService)
	
	// 認証ミドルウェアのモック
	authMiddleware := func(c *gin.Context) {
		// テスト用の認証情報を設定
		if userHeader := c.GetHeader("X-Test-User"); userHeader != "" {
			var user map[string]interface{}
			json.Unmarshal([]byte(userHeader), &user)
			c.Set("user", user)
		}
		c.Next()
	}
	
	r.GET("/api/relay", relayHandler.GetRelayRankings)
	r.POST("/api/relay", authMiddleware, relayHandler.RegisterRelayRankings)
	
	return r
}

func TestGetRelayRankings_Success(t *testing.T) {
	mockService := new(MockRelayService)
	router := setupRelayRouter(mockService)
	
	expectedRankings := map[int]int{
		1: 3, // 1位: 3年生
		2: 1, // 2位: 1年生
		3: 5, // 3位: 5年生
		4: 2, // 4位: 2年生
		5: 4, // 5位: 4年生
		6: 6, // 6位: 専・教
	}
	
	mockService.On("GetRelayRankings", "A").Return(expectedRankings, nil)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/relay?block=A", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response model.RelayRankResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "A", response.Block)
	assert.Equal(t, expectedRankings, response.Rankings)
	
	mockService.AssertExpectations(t)
}

func TestGetRelayRankings_MissingBlockParameter(t *testing.T) {
	mockService := new(MockRelayService)
	router := setupRelayRouter(mockService)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/relay", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "block parameter is required", response["error"])
}

func TestGetRelayRankings_InvalidBlock(t *testing.T) {
	mockService := new(MockRelayService)
	router := setupRelayRouter(mockService)
	
	mockService.On("GetRelayRankings", "C").Return(map[int]int{}, assert.AnError)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/relay?block=C", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	mockService.AssertExpectations(t)
}

func TestRegisterRelayRankings_Success_Superroot(t *testing.T) {
	mockService := new(MockRelayService)
	router := setupRelayRouter(mockService)
	
	rankings := map[int]int{
		1: 3, 2: 1, 3: 5, 4: 2, 5: 4, 6: 6,
	}
	
	mockService.On("RegisterRelayRankings", "A", rankings).Return(nil)
	
	requestBody := model.RelayRankRequest{Rankings: rankings}
	jsonBody, _ := json.Marshal(requestBody)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/relay?block=A", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	// superrootユーザーとして認証
	userInfo := map[string]interface{}{
		"role":           "superroot",
		"username":       "superroot",
		"assigned_sport": "",
	}
	userJSON, _ := json.Marshal(userInfo)
	req.Header.Set("X-Test-User", string(userJSON))
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "relay rankings registered successfully", response["message"])
	
	mockService.AssertExpectations(t)
}

func TestRegisterRelayRankings_Success_AdminRelay(t *testing.T) {
	mockService := new(MockRelayService)
	router := setupRelayRouter(mockService)
	
	rankings := map[int]int{
		1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6,
	}
	
	mockService.On("RegisterRelayRankings", "B", rankings).Return(nil)
	
	requestBody := model.RelayRankRequest{Rankings: rankings}
	jsonBody, _ := json.Marshal(requestBody)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/relay?block=B", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	// admin_relayユーザーとして認証
	userInfo := map[string]interface{}{
		"role":           "admin",
		"username":       "admin_relay",
		"assigned_sport": "relay",
	}
	userJSON, _ := json.Marshal(userInfo)
	req.Header.Set("X-Test-User", string(userJSON))
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	mockService.AssertExpectations(t)
}

func TestRegisterRelayRankings_Unauthorized_NoUser(t *testing.T) {
	mockService := new(MockRelayService)
	router := setupRelayRouter(mockService)
	
	rankings := map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6}
	requestBody := model.RelayRankRequest{Rankings: rankings}
	jsonBody, _ := json.Marshal(requestBody)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/relay?block=A", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRegisterRelayRankings_Forbidden_Student(t *testing.T) {
	mockService := new(MockRelayService)
	router := setupRelayRouter(mockService)
	
	rankings := map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6}
	requestBody := model.RelayRankRequest{Rankings: rankings}
	jsonBody, _ := json.Marshal(requestBody)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/relay?block=A", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	// studentユーザーとして認証
	userInfo := map[string]interface{}{
		"role":           "student",
		"username":       "student",
		"assigned_sport": "",
	}
	userJSON, _ := json.Marshal(userInfo)
	req.Header.Set("X-Test-User", string(userJSON))
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRegisterRelayRankings_Forbidden_WrongAdmin(t *testing.T) {
	mockService := new(MockRelayService)
	router := setupRelayRouter(mockService)
	
	rankings := map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6}
	requestBody := model.RelayRankRequest{Rankings: rankings}
	jsonBody, _ := json.Marshal(requestBody)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/relay?block=A", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	// 他のスポーツのadminユーザーとして認証
	userInfo := map[string]interface{}{
		"role":           "admin",
		"username":       "admin_volleyball",
		"assigned_sport": "volleyball",
	}
	userJSON, _ := json.Marshal(userInfo)
	req.Header.Set("X-Test-User", string(userJSON))
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRegisterRelayRankings_MissingBlockParameter(t *testing.T) {
	mockService := new(MockRelayService)
	router := setupRelayRouter(mockService)
	
	rankings := map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6}
	requestBody := model.RelayRankRequest{Rankings: rankings}
	jsonBody, _ := json.Marshal(requestBody)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/relay", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	userInfo := map[string]interface{}{
		"role":           "superroot",
		"username":       "superroot",
		"assigned_sport": "",
	}
	userJSON, _ := json.Marshal(userInfo)
	req.Header.Set("X-Test-User", string(userJSON))
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegisterRelayRankings_InvalidJSON(t *testing.T) {
	mockService := new(MockRelayService)
	router := setupRelayRouter(mockService)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/relay?block=A", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	
	userInfo := map[string]interface{}{
		"role":           "superroot",
		"username":       "superroot",
		"assigned_sport": "",
	}
	userJSON, _ := json.Marshal(userInfo)
	req.Header.Set("X-Test-User", string(userJSON))
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegisterRelayRankings_ServiceError(t *testing.T) {
	mockService := new(MockRelayService)
	router := setupRelayRouter(mockService)
	
	rankings := map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6}
	
	mockService.On("RegisterRelayRankings", "A", rankings).Return(assert.AnError)
	
	requestBody := model.RelayRankRequest{Rankings: rankings}
	jsonBody, _ := json.Marshal(requestBody)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/relay?block=A", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	userInfo := map[string]interface{}{
		"role":           "superroot",
		"username":       "superroot",
		"assigned_sport": "",
	}
	userJSON, _ := json.Marshal(userInfo)
	req.Header.Set("X-Test-User", string(userJSON))
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	mockService.AssertExpectations(t)
}