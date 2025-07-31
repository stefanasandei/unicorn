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
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// KeyVersion represents a key version for a user
type KeyVersion struct {
	ID        uuid.UUID  `json:"id" gorm:"type:text;primaryKey"`
	UserID    uuid.UUID  `json:"user_id" gorm:"type:text;not null;index"`
	Version   int        `json:"version" gorm:"not null"`
	KeyHash   string     `json:"key_hash" gorm:"not null;type:text"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	IsActive  bool       `json:"is_active" gorm:"not null;default:true"`
}

// KeyManager handles key rotation and versioning
type KeyManager struct {
	db    *gorm.DB
	mu    sync.RWMutex
	cache map[uuid.UUID]map[int][]byte // userID -> version -> key
}

// NewKeyManager creates a new key manager
func NewKeyManager(db *gorm.DB) (*KeyManager, error) {
	km := &KeyManager{
		db:    db,
		cache: make(map[uuid.UUID]map[int][]byte),
	}

	// Auto-migrate the key version table
	if err := db.AutoMigrate(&KeyVersion{}); err != nil {
		return nil, fmt.Errorf("failed to migrate key version schema: %w", err)
	}

	return km, nil
}

// deriveKey creates a 32-byte key from the user ID and version
func deriveKey(userID uuid.UUID, version int) []byte {
	keyData := fmt.Sprintf("%s:%d", userID.String(), version)
	h := sha256.Sum256([]byte(keyData))
	return h[:]
}

// GetOrCreateKey gets the key for a user and version, creating it if it doesn't exist
func (km *KeyManager) GetOrCreateKey(userID uuid.UUID, version int) ([]byte, error) {
	km.mu.RLock()
	if userKeys, exists := km.cache[userID]; exists {
		if key, exists := userKeys[version]; exists {
			km.mu.RUnlock()
			return key, nil
		}
	}
	km.mu.RUnlock()

	km.mu.Lock()
	defer km.mu.Unlock()

	// Double-check after acquiring write lock
	if userKeys, exists := km.cache[userID]; exists {
		if key, exists := userKeys[version]; exists {
			return key, nil
		}
	}

	// Check if key version exists in database
	var keyVersion KeyVersion
	err := km.db.Where("user_id = ? AND version = ?", userID, version).First(&keyVersion).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new key version
			key := deriveKey(userID, version)
			keyHash := fmt.Sprintf("%x", sha256.Sum256(key))

			keyVersion = KeyVersion{
				ID:        uuid.New(),
				UserID:    userID,
				Version:   version,
				KeyHash:   keyHash,
				CreatedAt: time.Now(),
				IsActive:  true,
			}

			if err := km.db.Create(&keyVersion).Error; err != nil {
				return nil, fmt.Errorf("failed to create key version: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to query key version: %w", err)
		}
	}

	// Generate the key
	key := deriveKey(userID, version)

	// Cache the key
	if km.cache[userID] == nil {
		km.cache[userID] = make(map[int][]byte)
	}
	km.cache[userID][version] = key

	return key, nil
}

// GetCurrentKeyVersion gets the current active key version for a user
func (km *KeyManager) GetCurrentKeyVersion(userID uuid.UUID) (int, error) {
	var keyVersion KeyVersion
	err := km.db.Where("user_id = ? AND is_active = ?", userID, true).
		Order("version DESC").
		First(&keyVersion).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create initial key version
			return km.createInitialKeyVersion(userID)
		}
		return 0, fmt.Errorf("failed to get current key version: %w", err)
	}

	return keyVersion.Version, nil
}

// createInitialKeyVersion creates the initial key version for a user
func (km *KeyManager) createInitialKeyVersion(userID uuid.UUID) (int, error) {
	key := deriveKey(userID, 1)
	keyHash := fmt.Sprintf("%x", sha256.Sum256(key))

	keyVersion := KeyVersion{
		ID:        uuid.New(),
		UserID:    userID,
		Version:   1,
		KeyHash:   keyHash,
		CreatedAt: time.Now(),
		IsActive:  true,
	}

	if err := km.db.Create(&keyVersion).Error; err != nil {
		return 0, fmt.Errorf("failed to create initial key version: %w", err)
	}

	// Cache the key
	if km.cache[userID] == nil {
		km.cache[userID] = make(map[int][]byte)
	}
	km.cache[userID][1] = key

	return 1, nil
}

// RotateKeys creates a new key version and re-encrypts all secrets
func (km *KeyManager) RotateKeys(userID uuid.UUID) error {
	// Get current version
	currentVersion, err := km.GetCurrentKeyVersion(userID)
	if err != nil {
		return fmt.Errorf("failed to get current key version: %w", err)
	}

	// Create new key version
	newVersion := currentVersion + 1
	newKey, err := km.GetOrCreateKey(userID, newVersion)
	if err != nil {
		return fmt.Errorf("failed to create new key: %w", err)
	}

	// Deactivate old key version
	if err := km.db.Model(&KeyVersion{}).
		Where("user_id = ? AND version = ?", userID, currentVersion).
		Update("is_active", false).Error; err != nil {
		return fmt.Errorf("failed to deactivate old key version: %w", err)
	}

	// Get all secrets for the user
	var secrets []struct {
		ID             uuid.UUID
		EncryptedValue string
		KeyVersion     int
	}

	if err := km.db.Table("secrets").
		Select("id, encrypted_value, key_version").
		Where("user_id = ?", userID).
		Find(&secrets).Error; err != nil {
		return fmt.Errorf("failed to get secrets for rotation: %w", err)
	}

	// Re-encrypt each secret with the new key
	for _, secret := range secrets {
		// Decrypt with old key
		oldKey, err := km.GetOrCreateKey(userID, secret.KeyVersion)
		if err != nil {
			return fmt.Errorf("failed to get old key for secret %s: %w", secret.ID, err)
		}

		plaintext, err := decryptSecret(secret.EncryptedValue, oldKey)
		if err != nil {
			return fmt.Errorf("failed to decrypt secret %s: %w", secret.ID, err)
		}

		// Encrypt with new key
		newEncrypted, err := encryptSecret(plaintext, newKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt secret %s with new key: %w", secret.ID, err)
		}

		// Update the secret
		if err := km.db.Model(&struct{ ID uuid.UUID }{}).
			Table("secrets").
			Where("id = ?", secret.ID).
			Updates(map[string]interface{}{
				"encrypted_value": newEncrypted,
				"key_version":     newVersion,
				"updated_at":      time.Now(),
			}).Error; err != nil {
			return fmt.Errorf("failed to update secret %s: %w", secret.ID, err)
		}
	}

	return nil
}

// GetKeyVersions gets all key versions for a user
func (km *KeyManager) GetKeyVersions(userID uuid.UUID) ([]KeyVersion, error) {
	var versions []KeyVersion
	err := km.db.Where("user_id = ?", userID).
		Order("version DESC").
		Find(&versions).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get key versions: %w", err)
	}

	return versions, nil
}

// encryptSecret encrypts a secret with the given key
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

// decryptSecret decrypts a secret with the given key
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
