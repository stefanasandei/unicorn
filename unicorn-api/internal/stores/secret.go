package stores

import (
	"fmt"
	"time"

	"unicorn-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SecretStoreInterface defines the interface for secret operations
type SecretStoreInterface interface {
	ListSecrets(userID uuid.UUID) ([]models.SecretResponse, error)
	CreateSecret(userID uuid.UUID, name, value, metadata string) (*models.Secret, error)
	GetSecret(userID, secretID uuid.UUID) (*models.Secret, string, error)
	UpdateSecret(userID, secretID uuid.UUID, value, metadata string) error
	DeleteSecret(userID, secretID uuid.UUID) error
	RotateKeys(userID uuid.UUID) error
	GetKeyVersions(userID uuid.UUID) ([]KeyVersion, error)
}

type SecretStore struct {
	db         *gorm.DB
	keyManager *KeyManager
}

// NewSecretStore creates a new SecretStore with SQLite
func NewSecretStore(dataSourceName string) (*SecretStore, error) {
	db, err := gorm.Open(sqlite.Open(dataSourceName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.AutoMigrate(&models.Secret{}); err != nil {
		return nil, fmt.Errorf("failed to migrate secret schema: %w", err)
	}

	keyManager, err := NewKeyManager(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create key manager: %w", err)
	}

	return &SecretStore{
		db:         db,
		keyManager: keyManager,
	}, nil
}

// CreateSecret stores a new encrypted secret
func (s *SecretStore) CreateSecret(userID uuid.UUID, name, value, metadata string) (*models.Secret, error) {
	// Get current key version
	keyVersion, err := s.keyManager.GetCurrentKeyVersion(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current key version: %w", err)
	}

	// Get key for current version
	key, err := s.keyManager.GetOrCreateKey(userID, keyVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	encrypted, err := encryptSecret(value, key)
	if err != nil {
		return nil, err
	}
	secret := &models.Secret{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Name:           name,
		EncryptedValue: encrypted,
		UserID:         userID,
		KeyVersion:     keyVersion,
		MetadataRaw:    metadata,
	}
	if err := s.db.Create(secret).Error; err != nil {
		return nil, err
	}
	return secret, nil
}

// GetSecret fetches and decrypts a secret for the owner
func (s *SecretStore) GetSecret(userID, secretID uuid.UUID) (*models.Secret, string, error) {
	var secret models.Secret
	if err := s.db.First(&secret, "id = ? AND user_id = ?", secretID, userID).Error; err != nil {
		return nil, "", err
	}

	// Get key for the secret's version
	key, err := s.keyManager.GetOrCreateKey(userID, secret.KeyVersion)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get key for version %d: %w", secret.KeyVersion, err)
	}

	plain, err := decryptSecret(secret.EncryptedValue, key)
	if err != nil {
		return nil, "", err
	}
	return &secret, plain, nil
}

// ListSecrets returns all secrets for a user (without values)
func (s *SecretStore) ListSecrets(userID uuid.UUID) ([]models.SecretResponse, error) {
	var secrets []models.Secret
	if err := s.db.Where("user_id = ?", userID).Find(&secrets).Error; err != nil {
		return nil, err
	}
	resp := make([]models.SecretResponse, len(secrets))
	for i, s := range secrets {
		resp[i] = models.SecretResponse{
			ID:         s.ID,
			Name:       s.Name,
			CreatedAt:  s.CreatedAt,
			UpdatedAt:  s.UpdatedAt,
			UserID:     s.UserID,
			KeyVersion: s.KeyVersion,
			Metadata:   s.MetadataRaw,
		}
	}
	return resp, nil
}

// UpdateSecret updates the value and/or metadata of a secret (owner only)
func (s *SecretStore) UpdateSecret(userID, secretID uuid.UUID, newValue, newMetadata string) error {
	var secret models.Secret
	if err := s.db.First(&secret, "id = ? AND user_id = ?", secretID, userID).Error; err != nil {
		return err
	}

	if newValue != "" {
		// Get current key version for new encryption
		currentKeyVersion, err := s.keyManager.GetCurrentKeyVersion(userID)
		if err != nil {
			return fmt.Errorf("failed to get current key version: %w", err)
		}

		key, err := s.keyManager.GetOrCreateKey(userID, currentKeyVersion)
		if err != nil {
			return fmt.Errorf("failed to get key: %w", err)
		}

		encrypted, err := encryptSecret(newValue, key)
		if err != nil {
			return err
		}
		secret.EncryptedValue = encrypted
		secret.KeyVersion = currentKeyVersion
	}
	if newMetadata != "" {
		secret.MetadataRaw = newMetadata
	}
	secret.UpdatedAt = time.Now()
	return s.db.Save(&secret).Error
}

// DeleteSecret removes a secret (owner only)
func (s *SecretStore) DeleteSecret(userID, secretID uuid.UUID) error {
	return s.db.Delete(&models.Secret{}, "id = ? AND user_id = ?", secretID, userID).Error
}

// RotateKeys rotates keys for a user
func (s *SecretStore) RotateKeys(userID uuid.UUID) error {
	return s.keyManager.RotateKeys(userID)
}

// GetKeyVersions gets all key versions for a user
func (s *SecretStore) GetKeyVersions(userID uuid.UUID) ([]KeyVersion, error) {
	return s.keyManager.GetKeyVersions(userID)
}
