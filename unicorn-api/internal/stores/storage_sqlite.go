package stores

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"unicorn-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GORMStorageStore implements storage operations using GORM for SQLite
// and writes file contents to disk.
type GORMStorageStore struct {
	db          *gorm.DB
	storagePath string
}

// NewGORMStorageStore creates a new GORMStorageStore
func NewGORMStorageStore(dataSourceName, storagePath string) (*GORMStorageStore, error) {
	db, err := gorm.Open(sqlite.Open(dataSourceName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database with GORM: %w", err)
	}

	err = db.AutoMigrate(&models.StorageBucket{}, &models.File{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate storage schema: %w", err)
	}

	return &GORMStorageStore{db: db, storagePath: storagePath}, nil
}

// CreateBucket creates a new storage bucket
func (s *GORMStorageStore) CreateBucket(bucket *models.StorageBucket) error {
	if bucket.ID == uuid.Nil {
		bucket.ID = uuid.New()
	}
	bucket.CreatedAt = time.Now()
	bucket.UpdatedAt = time.Now()
	result := s.db.Create(bucket)
	if result.Error != nil {
		return fmt.Errorf("failed to create bucket: %w", result.Error)
	}
	return nil
}

// GetBucketByID retrieves a bucket by its ID
func (s *GORMStorageStore) GetBucketByID(bucketID uuid.UUID) (*models.StorageBucket, error) {
	var bucket models.StorageBucket
	result := s.db.First(&bucket, "id = ?", bucketID)
	if result.Error != nil {
		return nil, fmt.Errorf("bucket not found: %w", result.Error)
	}
	return &bucket, nil
}

// ListFiles lists all files in a bucket
func (s *GORMStorageStore) ListFiles(bucketID uuid.UUID) ([]models.File, error) {
	var files []models.File
	result := s.db.Where("storage_bucket_id = ?", bucketID).Find(&files)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list files: %w", result.Error)
	}
	return files, nil
}

// SaveFile saves file metadata to DB and writes contents to disk
func (s *GORMStorageStore) SaveFile(bucketID uuid.UUID, fileHeader string, contentType string, fileData []byte) (*models.File, error) {
	fileID := uuid.New()
	filename := fileID.String() + "_" + fileHeader
	filePath := filepath.Join(s.storagePath, filename)
	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}
	fileModel := &models.File{
		ID:              fileID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Name:            fileHeader,
		Contents:        filename,
		Size:            int64(len(fileData)),
		ContentType:     contentType,
		StorageBucketID: bucketID,
	}
	if err := s.db.Create(fileModel).Error; err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save file metadata: %w", err)
	}
	return fileModel, nil
}

// GetFile retrieves file metadata by ID and bucket
func (s *GORMStorageStore) GetFile(bucketID, fileID uuid.UUID) (*models.File, error) {
	var file models.File
	result := s.db.First(&file, "id = ? AND storage_bucket_id = ?", fileID, bucketID)
	if result.Error != nil {
		return nil, fmt.Errorf("file not found: %w", result.Error)
	}
	return &file, nil
}

// DeleteFile deletes file metadata and removes file from disk
func (s *GORMStorageStore) DeleteFile(bucketID, fileID uuid.UUID) error {
	file, err := s.GetFile(bucketID, fileID)
	if err != nil {
		return err
	}
	filePath := filepath.Join(s.storagePath, file.Contents)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file from disk: %w", err)
	}
	if err := s.db.Delete(&models.File{}, "id = ? AND storage_bucket_id = ?", fileID, bucketID).Error; err != nil {
		return fmt.Errorf("failed to delete file metadata: %w", err)
	}
	return nil
}

// ListBucketsByUser returns all buckets owned by a user
func (s *GORMStorageStore) ListBucketsByUser(userID uuid.UUID) ([]models.StorageBucket, error) {
	var buckets []models.StorageBucket
	result := s.db.Where("user_id = ?", userID).Find(&buckets)
	if result.Error != nil {
		return nil, result.Error
	}
	return buckets, nil
}

// StoragePath returns the storage path
func (s *GORMStorageStore) StoragePath() string {
	return s.storagePath
}
