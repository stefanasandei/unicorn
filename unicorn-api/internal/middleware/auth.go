package middleware

import (
	"net/http"
	"strings"
	"unicorn-api/internal/auth"
	"unicorn-api/internal/config"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware returns a gin.HandlerFunc for JWT authentication
func AuthMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":       "Missing Authorization header",
				"status_code": http.StatusUnauthorized,
				"timestamp":   "2024-01-01T12:00:00Z",
			})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":       "Invalid Authorization header format. Expected 'Bearer ' followed by your token",
				"status_code": http.StatusUnauthorized,
				"timestamp":   "2024-01-01T12:00:00Z",
			})
			c.Abort()
			return
		}

		// Extract the token (remove "Bearer " prefix)
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":       "Missing token in Authorization header",
				"status_code": http.StatusUnauthorized,
				"timestamp":   "2024-01-01T12:00:00Z",
			})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := auth.ValidateToken(token, config)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":       "Invalid token",
				"status_code": http.StatusUnauthorized,
				"timestamp":   "2024-01-01T12:00:00Z",
			})
			c.Abort()
			return
		}

		// Store the claims in the context for use in handlers
		c.Set("claims", claims)
		c.Set("account_id", claims.AccountID)
		c.Set("role_id", claims.RoleID)

		c.Next()
	}
}

// GetClaimsFromContext is a helper function to extract claims from gin context
func GetClaimsFromContext(c *gin.Context) (*auth.Claims, bool) {
	claimsInterface, exists := c.Get("claims")
	if !exists {
		return nil, false
	}

	claims, ok := claimsInterface.(*auth.Claims)
	return claims, ok
}
