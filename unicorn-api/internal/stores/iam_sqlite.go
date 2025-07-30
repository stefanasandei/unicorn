package stores

import (
	"fmt"
	"unicorn-api/internal/auth"
	"unicorn-api/internal/config"
	"unicorn-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GORMIAMStore implements IAMStore using GORM for SQLite
type GORMIAMStore struct {
	db *gorm.DB
}

// NewGORMIAMStore creates a new GORMIAMStore
func NewGORMIAMStore(dataSourceName string) (*GORMIAMStore, error) {
	db, err := gorm.Open(sqlite.Open(dataSourceName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database with GORM: %w", err)
	}

	// AutoMigrate will create and update tables based on your models
	err = db.AutoMigrate(
		&models.Role{},
		&models.Organization{},
		&models.Account{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate database schema: %w", err)
	}

	return &GORMIAMStore{db: db}, nil
}

// CreateRole inserts a new role into the database
func (s *GORMIAMStore) CreateRole(role *models.Role) error {
	// GORM will automatically use the UUID in role.ID if it's set.
	// If ID is zero-value (empty UUID), GORM won't auto-generate it.
	// You might want to ensure ID is generated before calling Create.
	if role.ID == uuid.Nil {
		role.ID = uuid.New()
	}
	result := s.db.Create(role)
	if result.Error != nil {
		return fmt.Errorf("failed to create role: %w", result.Error)
	}
	return nil
}

// AssignRole assigns a role to an account (one-to-one relationship based on Account.RoleID)
func (s *GORMIAMStore) AssignRole(accountID, roleID string) error {
	accUUID, err := uuid.Parse(accountID)
	if err != nil {
		return fmt.Errorf("invalid account ID format: %w", err)
	}
	roleUUID, err := uuid.Parse(roleID)
	if err != nil {
		return fmt.Errorf("invalid role ID format: %w", err)
	}

	// Find the account
	var account models.Account
	result := s.db.Where("id = ?", accUUID).First(&account)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return models.ErrAccountNotFound // Use your custom error
		}
		return fmt.Errorf("failed to find account: %w", result.Error)
	}

	// Assign the role ID
	account.RoleID = roleUUID

	// Save the updated account
	result = s.db.Save(&account)
	if result.Error != nil {
		return fmt.Errorf("failed to assign role to account: %w", result.Error)
	}
	return nil
}

// GetRoleByName retrieves a role by its name
func (s *GORMIAMStore) GetRoleByName(name string) (*models.Role, error) {
	var role models.Role
	result := s.db.Where("name = ?", name).First(&role)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("role with name '%s' not found", name)
		}
		return nil, fmt.Errorf("failed to get role by name: %w", result.Error)
	}
	return &role, nil
}

// GetRoleByID retrieves a role by its ID
func (s *GORMIAMStore) GetRoleByID(roleID string) (*models.Role, error) {
	uid, err := uuid.Parse(roleID)
	if err != nil {
		return nil, fmt.Errorf("invalid role ID: %w", err)
	}
	var role models.Role
	result := s.db.Where("id = ?", uid).First(&role)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("role with ID '%s' not found", roleID)
		}
		return nil, fmt.Errorf("failed to get role by ID: %w", result.Error)
	}
	return &role, nil
}

// CreateOrganization inserts a new organization into the database
func (s *GORMIAMStore) CreateOrganization(org *models.Organization) error {
	if org.ID == uuid.Nil {
		org.ID = uuid.New()
	}
	result := s.db.Create(org)
	if result.Error != nil {
		return fmt.Errorf("failed to create organization: %w", result.Error)
	}
	return nil
}

// GetOrganizationByName retrieves an organization by its name
func (s *GORMIAMStore) GetOrganizationByName(name string) (
	*models.Organization,
	error,
) {
	var org models.Organization
	result := s.db.Where("name = ?", name).First(&org)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("organization with name '%s' not found", name)
		}
		return nil, fmt.Errorf("failed to get organization by name: %w", result.Error)
	}
	return &org, nil
}

// CreateAccount inserts a new account into the database
func (s *GORMIAMStore) CreateAccount(account *models.Account) error {
	// Check for duplicate email
	var existing models.Account
	if err := s.db.Where("email = ?", account.Email).First(&existing).Error; err == nil {
		return models.ErrDuplicateEmail
	}
	if account.ID == uuid.Nil {
		account.ID = uuid.New()
	}
	result := s.db.Create(account)
	if result.Error != nil {
		return fmt.Errorf("failed to create account: %w", result.Error)
	}
	return nil
}

// GetAccountByEmail retrieves an account by its email address
func (s *GORMIAMStore) GetAccountByEmail(email string) (*models.Account, error) {
	var account models.Account
	// Preload the Organization and Role if you need their data when fetching an account
	// This makes sure GORM fetches the associated data in one query or a few queries
	result := s.db.Where("email = ?", email).Preload("Organization").Preload("Role").First(&account)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, models.ErrAccountNotFound // Use your custom error
		}
		return nil, fmt.Errorf("failed to get account by email: %w", result.Error)
	}
	return &account, nil
}

// UpdateAccount updates an existing account in the database
func (s *GORMIAMStore) UpdateAccount(account *models.Account) error {
	result := s.db.Save(account) // Save will update if primary key exists, otherwise insert
	if result.Error != nil {
		return fmt.Errorf("failed to update account: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("account with ID '%s' not found for update", account.ID)
	}
	return nil
}

// --- New methods for GET /roles and /organizations ---

// GetAccountByID retrieves an account by its ID
func (s *GORMIAMStore) GetAccountByID(accountID string) (*models.Account, error) {
	uid, err := uuid.Parse(accountID)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID: %w", err)
	}
	var account models.Account
	result := s.db.Where("id = ?", uid).First(&account)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, models.ErrAccountNotFound
		}
		return nil, fmt.Errorf("failed to get account by ID: %w", result.Error)
	}
	return &account, nil
}

// GetRolesByOrganizationID returns all roles assigned to accounts in a given organization
func (s *GORMIAMStore) GetRolesByOrganizationID(orgID string) ([]models.Role, error) {
	uid, err := uuid.Parse(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}
	var accounts []models.Account
	if err := s.db.Where("organization_id = ?", uid).Find(&accounts).Error; err != nil {
		return nil, fmt.Errorf("failed to get accounts for organization: %w", err)
	}
	roleIDSet := make(map[uuid.UUID]struct{})
	for _, acc := range accounts {
		roleIDSet[acc.RoleID] = struct{}{}
	}
	if len(roleIDSet) == 0 {
		return []models.Role{}, nil
	}
	var roleIDs []uuid.UUID
	for rid := range roleIDSet {
		roleIDs = append(roleIDs, rid)
	}
	var roles []models.Role
	if err := s.db.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("failed to get roles by IDs: %w", err)
	}
	return roles, nil
}

// GetOrganizationByID retrieves an organization by its ID
func (s *GORMIAMStore) GetOrganizationByID(orgID string) (*models.Organization, error) {
	var org models.Organization
	result := s.db.Where("id = ?", orgID).First(&org)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, fmt.Errorf("failed to get organization by ID: %w", result.Error)
	}
	return &org, nil
}

// GetAccountsByOrganizationID returns all accounts for a given organization
func (s *GORMIAMStore) GetAccountsByOrganizationID(orgID string) ([]models.Account, error) {
	uid, err := uuid.Parse(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}
	var accounts []models.Account
	result := s.db.Where("organization_id = ?", uid).Find(&accounts)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get accounts by organization ID: %w", result.Error)
	}
	return accounts, nil
}

// SeedAdmin seeds the database with an admin organization, role, and user if not present
func (s *GORMIAMStore) SeedAdmin(cfg *config.Config) error {
	orgName := "admin-org"
	roleName := "admin"
	email := "admin@unicorn.local"
	password := "admin123"

	// 1. Create org if not exists
	org, err := s.GetOrganizationByName(orgName)
	if err != nil {
		org = &models.Organization{
			Name: orgName,
		}
		err = s.CreateOrganization(org)
		if err != nil {
			return err
		}
	}

	// 2. Create admin role if not exists
	role, err := s.GetRoleByName(roleName)
	if err != nil {
		role = &models.Role{
			Name:        roleName,
			Permissions: models.Permissions{models.Read, models.Write, models.Delete},
		}
		err = s.CreateRole(role)
		if err != nil {
			return err
		}
	}

	// 3. Create admin user if not exists
	acc, err := s.GetAccountByEmail(email)
	if err != nil {
		hash, _ := auth.HashPassword(password)
		acc = &models.Account{
			Name:           "Admin",
			Email:          email,
			Type:           models.AccountTypeUser,
			PasswordHash:   hash,
			OrganizationID: org.ID,
			RoleID:         role.ID,
		}
		err = s.CreateAccount(acc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *GORMIAMStore) DB() *gorm.DB {
	return s.db
}
