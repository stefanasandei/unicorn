package stores

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
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
}

type SecretStore struct {
	db *gorm.DB
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
	return &SecretStore{db: db}, nil
}

// deriveKey creates a 32-byte key from the user ID (for demo; use a real KMS in prod)
func deriveKey(userID uuid.UUID) []byte {
	h := sha256.Sum256([]byte(userID.String()))
	return h[:]
}

func encryptSecret(plainText string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func decryptSecret(cipherText string, key []byte) (string, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", errors.New("malformed ciphertext")
	}
	nonce, cipherData := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

// CreateSecret stores a new encrypted secret
func (s *SecretStore) CreateSecret(userID uuid.UUID, name, value, metadata string) (*models.Secret, error) {
	key := deriveKey(userID)
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
	key := deriveKey(userID)
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
			ID:        s.ID,
			Name:      s.Name,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
			UserID:    s.UserID,
			Metadata:  s.MetadataRaw,
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
	key := deriveKey(userID)
	if newValue != "" {
		encrypted, err := encryptSecret(newValue, key)
		if err != nil {
			return err
		}
		secret.EncryptedValue = encrypted
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
