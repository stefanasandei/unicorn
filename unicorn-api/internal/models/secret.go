package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Secret represents an encrypted secret stored in the database.
// swagger:model Secret
// @description An encrypted secret owned by a user.
type Secret struct {
	// The unique identifier of the secret
	ID uuid.UUID `json:"id" gorm:"type:text;primaryKey"`
	// The creation timestamp
	CreatedAt time.Time `json:"created_at"`
	// The last update timestamp
	UpdatedAt time.Time `json:"updated_at"`
	// The name/key of the secret
	Name string `json:"name" gorm:"not null;type:text"`
	// The encrypted value of the secret
	EncryptedValue string `json:"-" gorm:"not null;type:text"`
	// The ID of the user who owns the secret
	UserID uuid.UUID `json:"user_id" gorm:"type:text;not null;index"`
	// The key version used to encrypt this secret
	KeyVersion int `json:"key_version" gorm:"not null;default:1"`
	// Additional metadata for the secret (JSON)
	Metadata    map[string]string `gorm:"-" json:"metadata"`
	MetadataRaw string            `gorm:"type:text" json:"-"`
}

// SecretResponse represents the public view of a secret (without the encrypted value)
// swagger:model SecretResponse
// @description A secret response without sensitive data.
type SecretResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	UserID     uuid.UUID `json:"user_id"`
	KeyVersion int       `json:"key_version"`
	Metadata   string    `json:"metadata,omitempty"`
}

// CreateSecretRequest represents the request to create a new secret
// swagger:model CreateSecretRequest
type CreateSecretRequest struct {
	Name     string `json:"name" binding:"required"`
	Value    string `json:"value" binding:"required"`
	Metadata string `json:"metadata,omitempty"`
}

// UpdateSecretRequest represents the request to update a secret
// swagger:model UpdateSecretRequest
type UpdateSecretRequest struct {
	Name     string `json:"name,omitempty"`
	Value    string `json:"value,omitempty"`
	Metadata string `json:"metadata,omitempty"`
}

// SecretValueResponse represents the decrypted secret value response
// swagger:model SecretValueResponse
type SecretValueResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Value string    `json:"value"`
}

// SecretBodyRequest represents body required to create a secret
// swagger:model SecretBodyRequest
type SecretBodyRequest struct {
	Name     string `json:"name" binding:"required"`
	Value    string `json:"value" binding:"required"`
	Metadata string `json:"metadata"`
}

// UpdateSecretBody represents body required to update a secret
// swagger:model UpdateSecretBody
type UpdateSecretBody struct {
	Value    string `json:"value"`
	Metadata string `json:"metadata"`
}

// GORM hooks to marshal/unmarshal Metadata
func (s *Secret) BeforeSave(tx *gorm.DB) (err error) {
	if s.Metadata != nil {
		b, err := json.Marshal(s.Metadata)
		if err != nil {
			return err
		}
		s.MetadataRaw = string(b)
	}
	return nil
}

func (s *Secret) AfterFind(tx *gorm.DB) (err error) {
	if s.MetadataRaw != "" {
		var m map[string]string
		err := json.Unmarshal([]byte(s.MetadataRaw), &m)
		if err != nil {
			return err
		}
		s.Metadata = m
	}
	return nil
}
