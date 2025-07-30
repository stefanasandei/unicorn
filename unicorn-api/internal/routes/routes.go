package routes

import (
	"unicorn-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, iamHandler *handlers.IAMHandler, storageHandler *handlers.StorageHandler, computeHandler *handlers.ComputeHandler, lambdaHandler *handlers.LambdaHandler, secretHandler *handlers.SecretsHandler) {
	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// 1. IAM
		v1.POST("/roles", iamHandler.CreateRole)
		v1.GET("/roles", iamHandler.GetRoles)
		v1.POST("/roles/assign", iamHandler.AssignRole)
		v1.POST("/organizations", iamHandler.CreateOrganization)
		v1.GET("/organizations", iamHandler.GetOrganizations)
		v1.POST("/organizations/:org_id/users", iamHandler.CreateUserInOrg)
		v1.POST("/login", iamHandler.Login)
		v1.POST("/token/refresh", iamHandler.RefreshToken)
		v1.GET("/token/validate", iamHandler.ValidateToken)
		v1.GET("/debug/token", iamHandler.GetDebugToken)

		// 2. Secrets Manager
		v1.GET("/secrets", secretHandler.ListSecrets)
		v1.POST("/secrets", secretHandler.CreateSecret)
		v1.GET("/secrets/:id", secretHandler.ReadSecret)
		v1.PUT("/secrets/:id", secretHandler.UpdateSecret)
		v1.DELETE("/secrets/:id", secretHandler.DeleteSecret)

		// 3. Storage
		v1.GET("/buckets", storageHandler.ListBucketsHandler)
		v1.POST("/buckets", storageHandler.CreateBucketHandler)
		v1.POST("/buckets/:bucket_id/files", storageHandler.UploadFileHandler)
		v1.GET("/buckets/:bucket_id/files", storageHandler.ListFilesHandler)
		v1.GET("/buckets/:bucket_id/files/:file_id", storageHandler.DownloadFileHandler)
		v1.DELETE("/buckets/:bucket_id/files/:file_id", storageHandler.DeleteFileHandler)

		// 4. Compute
		v1.POST("/api/v1/compute/create", computeHandler.CreateCompute)
		v1.GET("/api/v1/compute/list", computeHandler.ListCompute)

		// 5. Lambda
		v1.POST("/lambda/execute", lambdaHandler.ExecuteLambda)
		v1.POST("/lambda/test", lambdaHandler.TestLambda)
	}
}
