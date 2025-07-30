package stores

import (
	"os"
	"testing"
	"time"
	"unicorn-api/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) (*GORMIAMStore, func()) {
	// Create a temporary database file
	dbPath := "test_" + uuid.New().String() + ".db"

	store, err := NewGORMIAMStore(dbPath)
	require.NoError(t, err)

	// Return cleanup function
	cleanup := func() {
		store.db.Migrator().DropTable(&models.Role{}, &models.Organization{}, &models.Account{})
		os.Remove(dbPath)
	}

	return store, cleanup
}

func TestCreateRole(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	tests := []struct {
		name        string
		role        *models.Role
		expectError bool
	}{
		{
			name: "valid role",
			role: &models.Role{
				ID:          uuid.New(),
				Name:        "admin",
				Permissions: models.Permissions{models.Read, models.Write, models.Delete},
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			expectError: false,
		},
		{
			name: "role with empty ID",
			role: &models.Role{
				Name:        "user",
				Permissions: models.Permissions{models.Read},
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			expectError: false, // Should auto-generate ID
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.CreateRole(tt.role)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, tt.role.ID)
			}
		})
	}
}

func TestGetRoleByName(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a test role
	role := &models.Role{
		ID:          uuid.New(),
		Name:        "test-role",
		Permissions: models.Permissions{models.Read, models.Write},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := store.CreateRole(role)
	require.NoError(t, err)

	tests := []struct {
		name        string
		roleName    string
		expectError bool
		expectNil   bool
	}{
		{
			name:        "existing role",
			roleName:    "test-role",
			expectError: false,
			expectNil:   false,
		},
		{
			name:        "non-existing role",
			roleName:    "non-existing",
			expectError: true,
			expectNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := store.GetRoleByName(tt.roleName)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if tt.expectNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tt.roleName, result.Name)
			}
		})
	}
}

func TestCreateOrganization(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	tests := []struct {
		name        string
		org         *models.Organization
		expectError bool
	}{
		{
			name: "valid organization",
			org: &models.Organization{
				ID:        uuid.New(),
				Name:      "Test Corp",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectError: false,
		},
		{
			name: "organization with empty ID",
			org: &models.Organization{
				Name:      "Another Corp",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectError: false, // Should auto-generate ID
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.CreateOrganization(tt.org)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, tt.org.ID)
			}
		})
	}
}

func TestGetOrganizationByName(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a test organization
	org := &models.Organization{
		ID:        uuid.New(),
		Name:      "test-org",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := store.CreateOrganization(org)
	require.NoError(t, err)

	tests := []struct {
		name        string
		orgName     string
		expectError bool
		expectNil   bool
	}{
		{
			name:        "existing organization",
			orgName:     "test-org",
			expectError: false,
			expectNil:   false,
		},
		{
			name:        "non-existing organization",
			orgName:     "non-existing",
			expectError: true,
			expectNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := store.GetOrganizationByName(tt.orgName)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if tt.expectNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tt.orgName, result.Name)
			}
		})
	}
}

func TestCreateAccount(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create required dependencies
	org := &models.Organization{
		ID:        uuid.New(),
		Name:      "test-org",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := store.CreateOrganization(org)
	require.NoError(t, err)

	role := &models.Role{
		ID:          uuid.New(),
		Name:        "test-role",
		Permissions: models.Permissions{models.Read},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = store.CreateRole(role)
	require.NoError(t, err)

	tests := []struct {
		name        string
		account     *models.Account
		expectError bool
	}{
		{
			name: "valid user account",
			account: &models.Account{
				ID:             uuid.New(),
				Name:           "John Doe",
				Email:          "john@example.com",
				Type:           models.AccountTypeUser,
				PasswordHash:   "hashed_password",
				OrganizationID: org.ID,
				RoleID:         role.ID,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			expectError: false,
		},
		{
			name: "valid bot account",
			account: &models.Account{
				ID:             uuid.New(),
				Name:           "Test Bot",
				Email:          "bot@example.com",
				Type:           models.AccountTypeBot,
				PasswordHash:   "bot_secret",
				OrganizationID: org.ID,
				RoleID:         role.ID,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.CreateAccount(tt.account)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, tt.account.ID)
			}
		})
	}
}

func TestGetAccountByEmail(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create required dependencies
	org := &models.Organization{
		ID:        uuid.New(),
		Name:      "test-org",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := store.CreateOrganization(org)
	require.NoError(t, err)

	role := &models.Role{
		ID:          uuid.New(),
		Name:        "test-role",
		Permissions: models.Permissions{models.Read},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = store.CreateRole(role)
	require.NoError(t, err)

	// Create a test account
	account := &models.Account{
		ID:             uuid.New(),
		Name:           "Test User",
		Email:          "test@example.com",
		Type:           models.AccountTypeUser,
		PasswordHash:   "hashed_password",
		OrganizationID: org.ID,
		RoleID:         role.ID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	err = store.CreateAccount(account)
	require.NoError(t, err)

	tests := []struct {
		name        string
		email       string
		expectError bool
		expectNil   bool
	}{
		{
			name:        "existing account",
			email:       "test@example.com",
			expectError: false,
			expectNil:   false,
		},
		{
			name:        "non-existing account",
			email:       "nonexistent@example.com",
			expectError: true,
			expectNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := store.GetAccountByEmail(tt.email)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if tt.expectNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tt.email, result.Email)
			}
		})
	}
}

func TestAssignRole(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create required dependencies
	org := &models.Organization{
		ID:        uuid.New(),
		Name:      "test-org",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := store.CreateOrganization(org)
	require.NoError(t, err)

	role1 := &models.Role{
		ID:          uuid.New(),
		Name:        "role1",
		Permissions: models.Permissions{models.Read},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = store.CreateRole(role1)
	require.NoError(t, err)

	role2 := &models.Role{
		ID:          uuid.New(),
		Name:        "role2",
		Permissions: models.Permissions{models.Read, models.Write},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = store.CreateRole(role2)
	require.NoError(t, err)

	// Create a test account
	account := &models.Account{
		ID:             uuid.New(),
		Name:           "Test User",
		Email:          "test@example.com",
		Type:           models.AccountTypeUser,
		PasswordHash:   "hashed_password",
		OrganizationID: org.ID,
		RoleID:         role1.ID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	err = store.CreateAccount(account)
	require.NoError(t, err)

	tests := []struct {
		name        string
		accountID   string
		roleID      string
		expectError bool
	}{
		{
			name:        "valid role assignment",
			accountID:   account.ID.String(),
			roleID:      role2.ID.String(),
			expectError: false,
		},
		{
			name:        "invalid account ID",
			accountID:   "invalid-uuid",
			roleID:      role2.ID.String(),
			expectError: true,
		},
		{
			name:        "invalid role ID",
			accountID:   account.ID.String(),
			roleID:      "invalid-uuid",
			expectError: true,
		},
		{
			name:        "non-existing account",
			accountID:   uuid.New().String(),
			roleID:      role2.ID.String(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.AssignRole(tt.accountID, tt.roleID)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify the role was actually assigned
				updatedAccount, err := store.GetAccountByEmail(account.Email)
				require.NoError(t, err)
				assert.Equal(t, tt.roleID, updatedAccount.RoleID.String())
			}
		})
	}
}

func TestUpdateAccount(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create required dependencies
	org := &models.Organization{
		ID:        uuid.New(),
		Name:      "test-org",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := store.CreateOrganization(org)
	require.NoError(t, err)

	role := &models.Role{
		ID:          uuid.New(),
		Name:        "test-role",
		Permissions: models.Permissions{models.Read},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = store.CreateRole(role)
	require.NoError(t, err)

	// Create a test account
	account := &models.Account{
		ID:             uuid.New(),
		Name:           "Original Name",
		Email:          "test@example.com",
		Type:           models.AccountTypeUser,
		PasswordHash:   "hashed_password",
		OrganizationID: org.ID,
		RoleID:         role.ID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	err = store.CreateAccount(account)
	require.NoError(t, err)

	// Test updating the account
	account.Name = "Updated Name"
	account.UpdatedAt = time.Now()

	err = store.UpdateAccount(account)
	assert.NoError(t, err)

	// Verify the update
	updatedAccount, err := store.GetAccountByEmail(account.Email)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", updatedAccount.Name)
}
