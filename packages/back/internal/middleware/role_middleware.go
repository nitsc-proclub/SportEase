package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuperrootOnly は superroot ロールを持つユーザーのみを許可するミドルウェアです。
// TableTennisAdminOnly は superroot または卓球担当のadminロールを持つユーザーのみを許可するミドルウェアです。
func TableTennisAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: Role not found"})
			return
		}

		userRole := role.(string)
		if userRole == "superroot" {
			c.Next()
			return
		}

		if userRole == "admin" {
			assignedSport, sportExists := c.Get("assigned_sport")
			if sportExists && assignedSport.(string) == "table_tennis" {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: You don't have access to this resource"})
	}
}

func SuperrootOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: Role not found"})
			return
		}

		userRole := role.(string)
		if userRole == "superroot" {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: You don't have access to this resource"})
	}
}
