package middleware

import (
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iamrichmon/subscription-api/internal/auth"
	"github.com/iamrichmon/subscription-api/internal/utils"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": utils.ErrInvalidAuthHeader.Error()})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": utils.ErrInvalidAuthHeader.Error()})
			c.Abort()
			return
		}
		tokenString := parts[1]

		claims, err := auth.ParseToken(tokenString, jwtSecret)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": utils.ErrInvalidCredentials.Error()})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
	}
}
