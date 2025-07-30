package handlers

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"unicorn-api/internal/auth"
	"unicorn-api/internal/config"
	"unicorn-api/internal/models"
	"unicorn-api/internal/stores"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StorageHandler struct {
	store   *stores.GORMStorageStore
	iamStore stores.IAMStore
	config  *config.Config
}

func NewStorageHandler(store *stores.GORMStorageStore, iamStore stores.IAMStore, config *config.Config) *StorageHandler {
	return &StorageHandler{store: store, iamStore: iamStore, config: config}
}

// Helper to extract claims from Authorization header
func (h *StorageHandler) getClaimsFromRequest(c *gin.Context) (*auth.Claims, error) {
	header := c.GetHeader("Authorization")
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		return nil, models.ErrTokenInvalid
	}
	token := strings.TrimPrefix(header, "Bearer ")
	return auth.ValidateToken(token, h.config)
}

// Helper to check permission
func hasPermission(permissions models.Permissions, required models.Permission) bool {
	for _, p := range permissions {
		if p == required {
			return true
		}
	}
	return false
}

// getUserPermissions fetches user permissions from their role
func (h *StorageHandler) getUserPermissions(userID uuid.UUID) (models.Permissions, error) {
	// For now, we'll use a simple approach - fetch the account and its role
	// In a real app, you might want to cache this
	account, err := h.iamStore.GetAccountByID(userID.String())
	if err != nil {
		return nil, err
	}
	
	role, err := h.iamStore.GetRoleByID(account.RoleID.String())
	if err != nil {
		return nil, err
	}
	
	return role.Permissions, nil
}

// CreateBucketRequest represents the request body for creating a bucket
// swagger:model
// @description Request to create a new storage bucket
// @name CreateBucketRequest
type CreateBucketRequest struct {
	// The name of the bucket
	// example: my-bucket
	Name string `json:"name" binding:"required"`
}

// ListBucketsHandler lists all buckets owned by the authenticated user
// @Summary List buckets
// @Description List all storage buckets owned by the authenticated user
// @Tags storage
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Success 200 {array} models.StorageBucket
// @Failure 401 {object} ErrorResponse
// @Router /buckets [get]
func (h *StorageHandler) ListBucketsHandler(c *gin.Context) {
	claims, err := h.getClaimsFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid or missing token"})
		return
	}
	userID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid user id in token"})
		return
	}
	buckets, err := h.store.ListBucketsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, buckets)
}

// CreateBucketHandler handles bucket creation
// @Summary Create bucket
// @Description Create a new storage bucket
// @Tags storage
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Param bucket body CreateBucketRequest true "Bucket name"
// @Success 201 {object} models.StorageBucket
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /buckets [post]
func (h *StorageHandler) CreateBucketHandler(c *gin.Context) {
	claims, err := h.getClaimsFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid or missing token"})
		return
	}
	
	userID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid user id in token"})
		return
	}
	
	// Check user permissions
	userPermissions, err := h.getUserPermissions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch user permissions"})
		return
	}
	
	if !hasPermission(userPermissions, models.Write) {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "permission denied: write access required"})
		return
	}
	
	var req CreateBucketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request"})
		return
	}
	
	bucket := models.StorageBucket{
		Name:   req.Name,
		UserID: userID,
		Files:  []models.File{},
	}
	if err := h.store.CreateBucket(&bucket); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, bucket)
}

// UploadFileHandler handles file uploads to a bucket
// @Summary Upload file
// @Description Upload a file to a storage bucket
// @Tags storage
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Param bucket_id path string true "Bucket ID"
// @Param file formData file true "File to upload"
// @Success 201 {object} models.File
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /buckets/{bucket_id}/files [post]
func (h *StorageHandler) UploadFileHandler(c *gin.Context) {
	claims, err := h.getClaimsFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid or missing token"})
		return
	}

	userID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid user id in token"})
		return
	}

	bucketID := c.Param("bucket_id")
	bucketUUID, err := uuid.Parse(bucketID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid bucket id"})
		return
	}
	
	bucket, err := h.store.GetBucketByID(bucketUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "bucket not found"})
		return
	}
	
	// Check if user owns the bucket
	if bucket.UserID != userID {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "permission denied: bucket does not belong to user"})
		return
	}
	
	// Check user permissions
	userPermissions, err := h.getUserPermissions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch user permissions"})
		return
	}
	
	if !hasPermission(userPermissions, models.Write) {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "permission denied: write access required"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "file required"})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "cannot open file"})
		return
	}
	defer file.Close()
	fileData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "cannot read file"})
		return
	}
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	fileModel, err := h.store.SaveFile(bucketUUID, fileHeader.Filename, contentType, fileData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, fileModel)
}

// DownloadFileHandler handles file downloads from a bucket
// @Summary Download file
// @Description Download a file from a storage bucket
// @Tags storage
// @Produce octet-stream
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Param bucket_id path string true "Bucket ID"
// @Param file_id path string true "File ID"
// @Success 200 {file} file
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /buckets/{bucket_id}/files/{file_id} [get]
func (h *StorageHandler) DownloadFileHandler(c *gin.Context) {
	claims, err := h.getClaimsFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid or missing token"})
		return
	}

	userID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid user id in token"})
		return
	}

	bucketID := c.Param("bucket_id")
	fileID := c.Param("file_id")
	bucketUUID, err := uuid.Parse(bucketID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid bucket id"})
		return
	}
	fileUUID, err := uuid.Parse(fileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid file id"})
		return
	}
	
	bucket, err := h.store.GetBucketByID(bucketUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "bucket not found"})
		return
	}
	
	// Check if user owns the bucket
	if bucket.UserID != userID {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "permission denied: bucket does not belong to user"})
		return
	}
	
	// Check user permissions
	userPermissions, err := h.getUserPermissions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch user permissions"})
		return
	}
	
	if !hasPermission(userPermissions, models.Read) {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "permission denied: read access required"})
		return
	}
	
	fileModel, err := h.store.GetFile(bucketUUID, fileUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "file not found"})
		return
	}
	filePath := filepath.Join(h.store.StoragePath(), fileModel.Contents)
	c.FileAttachment(filePath, fileModel.Name)
}

// ListFilesHandler lists files in a bucket
// @Summary List files
// @Description List all files in a storage bucket
// @Tags storage
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Param bucket_id path string true "Bucket ID"
// @Success 200 {array} models.File
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /buckets/{bucket_id}/files [get]
func (h *StorageHandler) ListFilesHandler(c *gin.Context) {
	claims, err := h.getClaimsFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid or missing token"})
		return
	}

	userID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid user id in token"})
		return
	}

	bucketID := c.Param("bucket_id")
	bucketUUID, err := uuid.Parse(bucketID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid bucket id"})
		return
	}
	
	bucket, err := h.store.GetBucketByID(bucketUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "bucket not found"})
		return
	}
	
	// Check if user owns the bucket
	if bucket.UserID != userID {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "permission denied: bucket does not belong to user"})
		return
	}
	
	// Check user permissions
	userPermissions, err := h.getUserPermissions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch user permissions"})
		return
	}
	
	if !hasPermission(userPermissions, models.Read) {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "permission denied: read access required"})
		return
	}
	
	files, err := h.store.ListFiles(bucketUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, files)
}

// DeleteFileHandler deletes a file from a bucket
// @Summary Delete file
// @Description Delete a file from a storage bucket
// @Tags storage
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Param bucket_id path string true "Bucket ID"
// @Param file_id path string true "File ID"
// @Success 204 {object} nil
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /buckets/{bucket_id}/files/{file_id} [delete]
func (h *StorageHandler) DeleteFileHandler(c *gin.Context) {
	claims, err := h.getClaimsFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid or missing token"})
		return
	}

	userID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid user id in token"})
		return
	}

	bucketID := c.Param("bucket_id")
	fileID := c.Param("file_id")
	bucketUUID, err := uuid.Parse(bucketID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid bucket id"})
		return
	}
	fileUUID, err := uuid.Parse(fileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid file id"})
		return
	}
	
	bucket, err := h.store.GetBucketByID(bucketUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "bucket not found"})
		return
	}
	
	// Check if user owns the bucket
	if bucket.UserID != userID {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "permission denied: bucket does not belong to user"})
		return
	}
	
	// Check user permissions
	userPermissions, err := h.getUserPermissions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch user permissions"})
		return
	}
	
	if !hasPermission(userPermissions, models.Delete) {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "permission denied: delete access required"})
		return
	}
	
	if err := h.store.DeleteFile(bucketUUID, fileUUID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
