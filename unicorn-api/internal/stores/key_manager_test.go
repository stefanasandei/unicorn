package stores

import (
	"testing"

	"unicorn-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupKeyManagerTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Auto-migrate the required tables
	if err := db.AutoMigrate(&KeyVersion{}, &models.Secret{}); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestKeyManager_CreateInitialKeyVersion(t *testing.T) {
	db := setupKeyManagerTestDB(t)
	km, err := NewKeyManager(db)
	if err != nil {
		t.Fatalf("Failed to create key manager: %v", err)
	}

	userID := uuid.New()
	version, err := km.GetCurrentKeyVersion(userID)
	if err != nil {
		t.Fatalf("Failed to get current key version: %v", err)
	}

	if version != 1 {
		t.Errorf("Expected initial version to be 1, got %d", version)
	}

	// Verify key version was created in database
	var keyVersion KeyVersion
	err = db.Where("user_id = ? AND version = ?", userID, 1).First(&keyVersion).Error
	if err != nil {
		t.Fatalf("Failed to find key version in database: %v", err)
	}

	if keyVersion.UserID != userID {
		t.Errorf("Expected user ID %s, got %s", userID, keyVersion.UserID)
	}

	if keyVersion.Version != 1 {
		t.Errorf("Expected version 1, got %d", keyVersion.Version)
	}

	if !keyVersion.IsActive {
		t.Error("Expected key version to be active")
	}
}

func TestKeyManager_KeyRotation(t *testing.T) {
	db := setupKeyManagerTestDB(t)
	km, err := NewKeyManager(db)
	if err != nil {
		t.Fatalf("Failed to create key manager: %v", err)
	}

	userID := uuid.New()

	// Get initial version
	initialVersion, err := km.GetCurrentKeyVersion(userID)
	if err != nil {
		t.Fatalf("Failed to get initial key version: %v", err)
	}

	if initialVersion != 1 {
		t.Errorf("Expected initial version to be 1, got %d", initialVersion)
	}

	// Rotate keys
	err = km.RotateKeys(userID)
	if err != nil {
		t.Fatalf("Failed to rotate keys: %v", err)
	}

	// Check new version
	newVersion, err := km.GetCurrentKeyVersion(userID)
	if err != nil {
		t.Fatalf("Failed to get new key version: %v", err)
	}

	if newVersion != 2 {
		t.Errorf("Expected new version to be 2, got %d", newVersion)
	}

	// Verify old version is deactivated
	var oldKeyVersion KeyVersion
	err = db.Where("user_id = ? AND version = ?", userID, 1).First(&oldKeyVersion).Error
	if err != nil {
		t.Fatalf("Failed to find old key version: %v", err)
	}

	if oldKeyVersion.IsActive {
		t.Error("Expected old key version to be inactive")
	}

	// Verify new version is active
	var newKeyVersion KeyVersion
	err = db.Where("user_id = ? AND version = ?", userID, 2).First(&newKeyVersion).Error
	if err != nil {
		t.Fatalf("Failed to find new key version: %v", err)
	}

	if !newKeyVersion.IsActive {
		t.Error("Expected new key version to be active")
	}
}

func TestKeyManager_GetKeyVersions(t *testing.T) {
	db := setupKeyManagerTestDB(t)
	km, err := NewKeyManager(db)
	if err != nil {
		t.Fatalf("Failed to create key manager: %v", err)
	}

	userID := uuid.New()

	// Create initial version
	_, err = km.GetCurrentKeyVersion(userID)
	if err != nil {
		t.Fatalf("Failed to get initial key version: %v", err)
	}

	// Rotate keys twice
	err = km.RotateKeys(userID)
	if err != nil {
		t.Fatalf("Failed to rotate keys first time: %v", err)
	}

	err = km.RotateKeys(userID)
	if err != nil {
		t.Fatalf("Failed to rotate keys second time: %v", err)
	}

	// Get all versions
	versions, err := km.GetKeyVersions(userID)
	if err != nil {
		t.Fatalf("Failed to get key versions: %v", err)
	}

	if len(versions) != 3 {
		t.Errorf("Expected 3 versions, got %d", len(versions))
	}

	// Check versions are in descending order
	for i, version := range versions {
		expectedVersion := 3 - i
		if version.Version != expectedVersion {
			t.Errorf("Expected version %d at index %d, got %d", expectedVersion, i, version.Version)
		}
	}

	// Check only the latest is active
	for i, version := range versions {
		if i == 0 && !version.IsActive {
			t.Error("Expected latest version to be active")
		}
		if i > 0 && version.IsActive {
			t.Errorf("Expected version %d to be inactive", version.Version)
		}
	}
}

func TestKeyManager_KeyCaching(t *testing.T) {
	db := setupKeyManagerTestDB(t)
	km, err := NewKeyManager(db)
	if err != nil {
		t.Fatalf("Failed to create key manager: %v", err)
	}

	userID := uuid.New()

	// Get key twice - should use cache on second call
	key1, err := km.GetOrCreateKey(userID, 1)
	if err != nil {
		t.Fatalf("Failed to get key first time: %v", err)
	}

	key2, err := km.GetOrCreateKey(userID, 1)
	if err != nil {
		t.Fatalf("Failed to get key second time: %v", err)
	}

	// Keys should be the same
	if string(key1) != string(key2) {
		t.Error("Expected cached keys to be identical")
	}
}
