package services

import (
	"encoding/json"

	"unicorn-api/internal/common/errors"
	"unicorn-api/internal/common/validation"
	"unicorn-api/internal/models"
	"unicorn-api/internal/stores"

	"github.com/google/uuid"
)

// SecretService handles secret-related business logic
type SecretService struct {
	store     stores.SecretStoreInterface
	validator *validation.Validator
}

// NewSecretService creates a new secret service
func NewSecretService(store stores.SecretStoreInterface) *SecretService {
	return &SecretService{
		store:     store,
		validator: validation.NewValidator(),
	}
}

// CreateSecret creates a new secret with proper validation
func (s *SecretService) CreateSecret(userID uuid.UUID, name, value, metadata string) (*models.Secret, error) {
	// Validate inputs
	if err := s.validator.ValidateSecretName(name); err != nil {
		return nil, err
	}

	if err := s.validator.ValidateSecretValue(value); err != nil {
		return nil, err
	}

	// Validate metadata JSON if provided
	if metadata != "" {
		if err := s.validator.ValidateJSON(metadata); err != nil {
			return nil, errors.ErrBadRequest.WithDetails("Invalid metadata JSON format")
		}
	}

	// Create secret
	secret, err := s.store.CreateSecret(userID, name, value, metadata)
	if err != nil {
		return nil, errors.ErrInternalError.WithDetails("Failed to create secret: " + err.Error())
	}

	return secret, nil
}

// GetSecret retrieves a secret with proper validation
func (s *SecretService) GetSecret(userID, secretID uuid.UUID) (*models.Secret, string, error) {
	secret, value, err := s.store.GetSecret(userID, secretID)
	if err != nil {
		return nil, "", errors.ErrResourceNotFound.WithDetails("Secret not found")
	}

	return secret, value, nil
}

// ListSecrets retrieves all secrets for a user
func (s *SecretService) ListSecrets(userID uuid.UUID) ([]models.SecretResponse, error) {
	secrets, err := s.store.ListSecrets(userID)
	if err != nil {
		return nil, errors.ErrInternalError.WithDetails("Failed to list secrets: " + err.Error())
	}

	return secrets, nil
}

// UpdateSecret updates a secret with proper validation
func (s *SecretService) UpdateSecret(userID, secretID uuid.UUID, value, metadata string) error {
	// Validate inputs
	if value != "" {
		if err := s.validator.ValidateSecretValue(value); err != nil {
			return err
		}
	}

	// Validate metadata JSON if provided
	if metadata != "" {
		if err := s.validator.ValidateJSON(metadata); err != nil {
			return errors.ErrBadRequest.WithDetails("Invalid metadata JSON format")
		}
	}

	// Check if secret exists
	_, _, err := s.store.GetSecret(userID, secretID)
	if err != nil {
		return errors.ErrResourceNotFound.WithDetails("Secret not found")
	}

	// Update secret
	if err := s.store.UpdateSecret(userID, secretID, value, metadata); err != nil {
		return errors.ErrInternalError.WithDetails("Failed to update secret: " + err.Error())
	}

	return nil
}

// DeleteSecret deletes a secret
func (s *SecretService) DeleteSecret(userID, secretID uuid.UUID) error {
	// Check if secret exists
	_, _, err := s.store.GetSecret(userID, secretID)
	if err != nil {
		return errors.ErrResourceNotFound.WithDetails("Secret not found")
	}

	if err := s.store.DeleteSecret(userID, secretID); err != nil {
		return errors.ErrInternalError.WithDetails("Failed to delete secret: " + err.Error())
	}

	return nil
}

// ValidateSecretMetadata validates and parses secret metadata
func (s *SecretService) ValidateSecretMetadata(metadata string) (map[string]string, error) {
	if metadata == "" {
		return nil, nil
	}

	var metadataMap map[string]string
	if err := json.Unmarshal([]byte(metadata), &metadataMap); err != nil {
		return nil, errors.ErrBadRequest.WithDetails("Invalid metadata JSON format")
	}

	return metadataMap, nil
}

// RotateKeys rotates all keys for a user
func (s *SecretService) RotateKeys(userID uuid.UUID) error {
	if err := s.store.RotateKeys(userID); err != nil {
		return errors.ErrInternalError.WithDetails("Failed to rotate keys: " + err.Error())
	}
	return nil
}

// GetKeyVersions gets all key versions for a user
func (s *SecretService) GetKeyVersions(userID uuid.UUID) ([]stores.KeyVersion, error) {
	versions, err := s.store.GetKeyVersions(userID)
	if err != nil {
		return nil, errors.ErrInternalError.WithDetails("Failed to get key versions: " + err.Error())
	}
	return versions, nil
}
