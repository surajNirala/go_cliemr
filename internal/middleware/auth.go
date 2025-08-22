// middleware/auth.go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/surajNirala/go_cliemr/pkg/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.RespondError(c, "error", http.StatusUnauthorized, "Authorization header missing", "Authorization header is required")
			c.Abort()
			return
		}

		// Expect "Bearer <token>"
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			utils.RespondError(c, "error", http.StatusUnauthorized, "Invalid token format", "Token must be in 'Bearer <token>' format")
			c.Abort()
			return
		}

		claims, err := utils.ValidateAccessToken(token)
		if err != nil {
			utils.RespondError(c, "error", http.StatusUnauthorized, "Invalid or expired token", err.Error())
			c.Abort()
			return
		}

		// Store claims in context (e.g., userID, role)
		c.Set("userID", claims.UserID)
		c.Set("flag", claims.Flag)
		c.Set("email", claims.Email)
		c.Next()
	}
}
