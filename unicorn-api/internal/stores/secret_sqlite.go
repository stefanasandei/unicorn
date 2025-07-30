package stores

import (
	"time"
	"unicorn-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GORMSecretsStore struct {
	db *gorm.DB
}

func NewGORMSecretsStore(dataSourceName string) (*GORMSecretsStore, error) {
	db, err := gorm.Open(sqlite.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&models.Secret{}); err != nil {
		return nil, err
	}
	return &GORMSecretsStore{db: db}, nil
}

type SecretsStore interface {
	ListSecrets(userID uuid.UUID) ([]models.Secret, error)
	CreateSecret(userID uuid.UUID, name, value string, metadata string) (*models.Secret, error)
	GetSecret(userID, secretID uuid.UUID) (*models.Secret, string, error)
	UpdateSecret(userID, secretID uuid.UUID, value string, metadata string) error
	DeleteSecret(userID, secretID uuid.UUID) error
}

// Example method signatures (implementations should be filled in as needed)
func (s *GORMSecretsStore) ListSecrets(userID uuid.UUID) ([]models.Secret, error) {
	var secrets []models.Secret
	err := s.db.Where("user_id = ?", userID).Find(&secrets).Error
	return secrets, err
}

func (s *GORMSecretsStore) CreateSecret(userID uuid.UUID, name, value string, metadata string) (*models.Secret, error) {
	secret := &models.Secret{
		ID:             uuid.New(),
		UserID:         userID,
		Name:           name,
		EncryptedValue: value, // For compatibility, but should be encrypted
		MetadataRaw:    metadata,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	err := s.db.Create(secret).Error
	return secret, err
}

func (s *GORMSecretsStore) GetSecret(userID, secretID uuid.UUID) (*models.Secret, string, error) {
	var secret models.Secret
	err := s.db.Where("id = ? AND user_id = ?", secretID, userID).First(&secret).Error
	return &secret, secret.EncryptedValue, err
}

func (s *GORMSecretsStore) UpdateSecret(userID, secretID uuid.UUID, value string, metadata string) error {
	return s.db.Model(&models.Secret{}).Where("id = ? AND user_id = ?", secretID, userID).Updates(map[string]interface{}{"encrypted_value": value, "metadata_raw": metadata, "updated_at": time.Now()}).Error
}

func (s *GORMSecretsStore) DeleteSecret(userID, secretID uuid.UUID) error {
	return s.db.Where("id = ? AND user_id = ?", secretID, userID).Delete(&models.Secret{}).Error
}
