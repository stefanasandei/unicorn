package routes

import (
	"unicorn-api/internal/config"
	"unicorn-api/internal/handlers"
	"unicorn-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, iamHandler *handlers.IAMHandler, storageHandler *handlers.StorageHandler, computeHandler *handlers.ComputeHandler, lambdaHandler *handlers.LambdaHandler, secretHandler *handlers.SecretsHandler, rdbHandler *handlers.RDBHandler, config *config.Config) {
	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		{
			v1.POST("/login", iamHandler.Login)
			v1.POST("/token/refresh", iamHandler.RefreshToken)
			v1.GET("/token/validate", iamHandler.ValidateToken)
			v1.GET("/debug/token", iamHandler.GetDebugToken)

			// Setup routes (no authentication required for initial setup)
			v1.POST("/organizations", iamHandler.CreateOrganization)
			v1.POST("/roles", iamHandler.CreateRole)
			v1.POST("/organizations/:org_id/users", iamHandler.CreateUserInOrg)
		}

		// Protected routes (authentication required)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(config))
		{
			// IAM routes
			protected.GET("/roles", iamHandler.GetRoles)
			protected.POST("/roles/assign", iamHandler.AssignRole)
			protected.GET("/organizations", iamHandler.GetOrganizations)

			// Secrets Manager routes
			protected.GET("/secrets", secretHandler.ListSecrets)
			protected.POST("/secrets", secretHandler.CreateSecret)
			protected.GET("/secrets/:id", secretHandler.ReadSecret)
			protected.PUT("/secrets/:id", secretHandler.UpdateSecret)
			protected.DELETE("/secrets/:id", secretHandler.DeleteSecret)

			// Storage routes
			protected.GET("/buckets", storageHandler.ListBucketsHandler)
			protected.POST("/buckets", storageHandler.CreateBucketHandler)
			protected.POST("/buckets/:bucket_id/files", storageHandler.UploadFileHandler)
			protected.GET("/buckets/:bucket_id/files", storageHandler.ListFilesHandler)
			protected.GET("/buckets/:bucket_id/files/:file_id", storageHandler.DownloadFileHandler)
			protected.DELETE("/buckets/:bucket_id/files/:file_id", storageHandler.DeleteFileHandler)

			// Compute routes
			protected.POST("/compute/create", computeHandler.CreateCompute)
			protected.GET("/compute/list", computeHandler.ListCompute)
			protected.DELETE("/compute/:id", computeHandler.DeleteCompute)

			// Lambda routes
			protected.POST("/lambda/execute", lambdaHandler.ExecuteLambda)
			protected.POST("/lambda/test", lambdaHandler.TestLambda)

			// RDB routes
			protected.POST("/rdb/create", rdbHandler.CreateRDB)
			protected.GET("/rdb/list", rdbHandler.ListRDB)
			protected.DELETE("/rdb/:id", rdbHandler.DeleteRDB)
		}
	}
}
