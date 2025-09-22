package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/saku0512/GYOUJI_HP/backend/internal/service" // serviceパッケージのClaimsを参照
)

// AuthMiddleware はJWTを検証するミドルウェアです。
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization cookie is required"})
			return
		}

		// トークンをパース・検証
		claims := &service.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Ginのコンテキストにユーザー情報を保存
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		// userが未セットならまとめてセット
		if _, exists := c.Get("user"); !exists {
			c.Set("user", map[string]interface{}{
				"user_id":        claims.UserID,
				"username":       claims.Username,
				"role":           claims.Role,
				"assigned_sport": claims.AssignedSport,
			})
		}

		// 次の処理へ進む
		c.Next()
	}
}
