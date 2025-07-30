package routes

import (
	"unicorn-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, iamHandler *handlers.IAMHandler) {
	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Example endpoints
		v1.GET("/hello", handlers.Hello)

		// 1. IAM
		v1.POST("/roles", iamHandler.CreateRole)
		v1.GET("/roles", iamHandler.GetRoles) // Added GET /roles
		v1.POST("/roles/assign", iamHandler.AssignRole)
		v1.POST("/organizations", iamHandler.CreateOrganization)
		v1.GET("/organizations", iamHandler.GetOrganizations) // Added GET /organizations
		v1.POST("/organizations/:org_id/users", iamHandler.CreateUserInOrg)
		v1.POST("/login", iamHandler.Login)
		v1.POST("/token/refresh", iamHandler.RefreshToken)
		v1.GET("/token/validate", iamHandler.ValidateToken)
	}
}
