package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/surajNirala/go_cliemr/pkg/utils"
)

// RequireRoles ensures only users with allowed roles can access route
func RequireRoles(allowedRoles ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Assume you set user flag from JWT in context
		flagValue, exists := c.Get("flag")
		if !exists {
			utils.RespondError(c, "error", http.StatusUnauthorized, "User role not found", "User role is required")
			c.Abort()
			return
		}

		var userFlag int
		switch v := flagValue.(type) {
		case int:
			userFlag = v
		case float64:
			userFlag = int(v) // if stored from JWT claims
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid flag type"})
			c.Abort()
			return
		}

		// Check if userFlag is in allowedRoles
		allowed := false
		for _, role := range allowedRoles {
			if role == userFlag {
				allowed = true
				break
			}
		}

		if !allowed {
			utils.RespondError(c, "error", http.StatusForbidden, "You are not authorized to access this resource", "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}
