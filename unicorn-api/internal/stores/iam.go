package stores

import (
	"unicorn-api/internal/models"
)

// IAMStore abstracts DB operations for IAM
// Implement this interface for the DB layer
type IAMStore interface {
	// Role management
	CreateRole(role *models.Role) error
	AssignRole(accountID, roleID string) error
	GetRoleByName(name string) (*models.Role, error)

	// Organization management
	CreateOrganization(org *models.Organization) error
	GetOrganizationByName(name string) (*models.Organization, error)

	// User management
	CreateAccount(account *models.Account) error
	GetAccountByEmail(email string) (*models.Account, error)
	UpdateAccount(account *models.Account) error

	// --- Added for GET /roles and /organizations ---
	GetAccountByID(accountID string) (*models.Account, error)
	GetRolesByOrganizationID(orgID string) ([]models.Role, error)
	GetOrganizationByID(orgID string) (*models.Organization, error)
	GetAccountsByOrganizationID(orgID string) ([]models.Account, error)
}
