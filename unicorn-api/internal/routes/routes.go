package routes

import (
	"unicorn-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, iamHandler *handlers.IAMHandler, storageHandler *handlers.StorageHandler) {
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

		// 2. Storage
		v1.GET("/buckets", storageHandler.ListBucketsHandler)
		v1.POST("/buckets", storageHandler.CreateBucketHandler)
		v1.POST("/buckets/:bucket_id/files", storageHandler.UploadFileHandler)
		v1.GET("/buckets/:bucket_id/files", storageHandler.ListFilesHandler)
		v1.GET("/buckets/:bucket_id/files/:file_id", storageHandler.DownloadFileHandler)
		v1.DELETE("/buckets/:bucket_id/files/:file_id", storageHandler.DeleteFileHandler)
	}
}
