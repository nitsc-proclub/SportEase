package handler

import (
	"net/http"

	"github.com/saku0512/GYOUJI_HP/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// UserHandler はユーザー関連のHTTPリクエストを処理します。
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler は UserHandler の新しいインスタンスを生成します。
func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{userService: s}
}

// LoginRequest はログインリクエストのボディを表します。
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login はログイン処理を行い、JWTを返します。
// @Summary User login
// @Description Authenticate a user and return a JWT token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Login Credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	token, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// JWTをHttpOnly Cookieに設定
	// localhostでの開発環境と本番環境の両方で動作するように設定
	// 本番環境では secure を true に、samesite を none にする必要がある
	secure := c.Request.Header.Get("X-Forwarded-Proto") == "https"
	c.SetCookie("jwt", token, 3600*24, "/", c.Request.Host, secure, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// Logout godoc
// @Summary Logout a user
// @Description Invalidate the user's session.
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/auth/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	// Cookieをクリア
	secure := c.Request.Header.Get("X-Forwarded-Proto") == "https"
	c.SetCookie("jwt", "", -1, "/", c.Request.Host, secure, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// GET /api/auth/me
// GetMe godoc
// @Summary Get UserData
// @Description Get UserData by AuthMiddleware
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/auth/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	user, exists := c.Get("user") // AuthMiddlewareでセットされている前提
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	c.JSON(http.StatusOK, user)
}
