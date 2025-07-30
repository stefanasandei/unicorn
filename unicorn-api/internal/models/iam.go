package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// 1. general errors for the IAM module
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountNotFound    = errors.New("account not found")
	ErrTokenInvalid       = errors.New("invalid token")
	ErrPermissionDenied   = errors.New("permission denied")
)

// 2. roles & permissions

// Permission represents the level of access a role has
// 0 = Read, 1 = Write, 2 = Delete
type Permission int

const (
	Read   Permission = iota // Read permission (0)
	Write                    // Write permission (1)
	Delete                   // Delete permission (2)
)

// String returns the string representation of the permission
func (p Permission) String() string {
	switch p {
	case Read:
		return "read"
	case Write:
		return "write"
	case Delete:
		return "delete"
	default:
		return "unknown"
	}
}

// Permissions represents a slice of Permission that can be marshalled/unmarshalled to JSON
type Permissions []Permission

// Value implements the driver.Valuer interface for database serialization.
func (p Permissions) Value() (driver.Value, error) {
	if len(p) == 0 {
		return "[]", nil // Store empty array as "[]"
	}
	j, err := json.Marshal(p)
	return string(j), err
}

// Scan implements the sql.Scanner interface for database deserialization.
func (p *Permissions) Scan(src interface{}) error {
	var source []byte
	switch src := src.(type) {
	case string:
		source = []byte(src)
	case []byte:
		source = src
	default:
		return errors.New("incompatible type for Permissions")
	}

	if len(source) == 0 {
		*p = []Permission{} // Handle empty string as empty slice
		return nil
	}

	return json.Unmarshal(source, p)
}

// Role represents a user role in the system with associated permissions
// swagger:model
type Role struct {
	// The unique ID of the role
	// example: 123e4567-e89b-12d3-a456-426614174000
	ID uuid.UUID `json:"id" gorm:"type:text;primaryKey"`
	// The creation timestamp
	// example: 2024-01-01T12:00:00Z
	CreatedAt time.Time `json:"created_at"`
	// The last update timestamp
	// example: 2024-01-01T12:00:00Z
	UpdatedAt time.Time `json:"updated_at"`
	// The name of the role (e.g., admin, user, moderator)
	// example: admin
	Name string `json:"name" gorm:"unique;not null;type:text"`
	// The permissions assigned to the role (0=Read, 1=Write, 2=Delete)
	// example: [0,1,2]
	Permissions Permissions `json:"permissions" gorm:"type:json"`
}

// 3. user & bots can have accounts

// AccountType represents the type of account in the system
type AccountType string

const (
	AccountTypeUser AccountType = "user" // Human user account
	AccountTypeBot  AccountType = "bot"  // Automated service account
)

// Account represents a user or bot account in the system
// swagger:model
type Account struct {
	// The unique ID of the account
	// example: 123e4567-e89b-12d3-a456-426614174000
	ID uuid.UUID `json:"id" gorm:"type:text;primaryKey"`
	// The creation timestamp
	// example: 2024-01-01T12:00:00Z
	CreatedAt time.Time `json:"created_at"`
	// The last update timestamp
	// example: 2024-01-01T12:00:00Z
	UpdatedAt time.Time `json:"updated_at"`

	// The display name of the account
	// example: John Doe
	Name string `json:"name" gorm:"type:text"`
	// The email address (required for user accounts, unique)
	// example: john.doe@example.com
	Email string `json:"email,omitempty" gorm:"unique;type:text"`
	// The type of account (user or bot)
	// example: user
	Type AccountType `json:"type" gorm:"type:text"`

	// Hashed and salted password for user accounts.
	// For service accounts, this might be empty or used for an initial secret.
	// Omitted from JSON output for security
	PasswordHash string `json:"-" gorm:"not null;type:text"`

	// Token related (JWT) - last successful login timestamp
	// example: 2024-01-01T12:00:00Z
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`

	// Foreign Keys
	// The ID of the organization this account belongs to
	// example: 123e4567-e89b-12d3-a456-426614174000
	OrganizationID uuid.UUID `json:"organization_id" gorm:"type:text"`
	// The ID of the role assigned to this account
	// example: 123e4567-e89b-12d3-a456-426614174000
	RoleID uuid.UUID `json:"role_id" gorm:"type:text"`

	// GORM Associations
	// The organization this account belongs to
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
	// The role assigned to this account
	Role Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// Organization represents a group or company that can contain multiple accounts
// swagger:model
type Organization struct {
	// The unique ID of the organization
	// example: 123e4567-e89b-12d3-a456-426614174000
	ID uuid.UUID `json:"id" gorm:"type:text;primaryKey"`
	// The creation timestamp
	// example: 2024-01-01T12:00:00Z
	CreatedAt time.Time `json:"created_at"`
	// The last update timestamp
	// example: 2024-01-01T12:00:00Z
	UpdatedAt time.Time `json:"updated_at"`

	// The name of the organization (unique)
	// example: Acme Corporation
	Name string `json:"name" gorm:"unique;not null;type:text"`

	// An organization can have many accounts
	// example: [{"id":"123e4567-e89b-12d3-a456-426614174000","name":"John Doe","email":"john@acme.com","type":"user"}]
	Accounts []Account `gorm:"foreignKey:OrganizationID" json:"accounts,omitempty"`
}
