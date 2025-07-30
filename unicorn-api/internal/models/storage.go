package models

import (
	"time"

	"github.com/google/uuid"
)

// StorageBucket represents a bucket for storing files.
// swagger:model StorageBucket
// @description A storage bucket owned by a user, containing files.
type StorageBucket struct {
	// The unique identifier of the bucket
	ID uuid.UUID `json:"id" gorm:"type:text;primaryKey"`
	// The creation timestamp
	CreatedAt time.Time `json:"created_at"`
	// The last update timestamp
	UpdatedAt time.Time `json:"updated_at"`
	// The name of the bucket
	Name string `json:"name" gorm:"unique;not null;type:text"`
	// The ID of the user who owns the bucket
	UserID uuid.UUID `json:"user_id" gorm:"type:text;not null;index"`
	// The files in the bucket
	Files []File `json:"files" gorm:"foreignKey:StorageBucketID"`
}

// File represents a file stored in a bucket.
// swagger:model File
// @description A file stored in a storage bucket.
type File struct {
	// The unique identifier of the file
	ID uuid.UUID `json:"id" gorm:"type:text;primaryKey"`
	// The creation timestamp
	CreatedAt time.Time `json:"created_at"`
	// The last update timestamp
	UpdatedAt time.Time `json:"updated_at"`
	// The name of the file
	Name string `json:"name" gorm:"not null;type:text"`
	// The file contents (base64 or text)
	Contents string `json:"contents" gorm:"not null;type:text"`
	// The size of the file in bytes
	Size int64 `json:"size" gorm:"not null"`
	// The MIME type of the file
	ContentType string `json:"content_type" gorm:"not null;type:text"`
	// The ID of the bucket this file belongs to
	StorageBucketID uuid.UUID `json:"storage_bucket_id" gorm:"type:text;not null;index"`
}
