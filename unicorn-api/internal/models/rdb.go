package models

// RDBType defines the type of database
type RDBType string

const (
	RDBTypePostgreSQL RDBType = "postgresql"
	RDBTypeMySQL      RDBType = "mysql"
)

// RDBPreset defines hardware constraints for RDB service
type RDBPreset string

const (
	RDBPresetMicro  RDBPreset = "micro"
	RDBPresetSmall  RDBPreset = "small"
	RDBPresetMedium RDBPreset = "medium"
)

// RDBVolume represents a volume configuration for the database
type RDBVolume struct {
	Name string `json:"name" binding:"required"`
	Size int    `json:"size" binding:"required,min=1"` // Size in MB
}

// RDBCreateRequest is the request body for creating an RDB instance
type RDBCreateRequest struct {
	Name        string            `json:"name,omitempty"`
	Type        RDBType           `json:"type" binding:"required"`
	Preset      RDBPreset         `json:"preset,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	Volumes     []RDBVolume       `json:"volumes,omitempty"`
	Port        string            `json:"port,omitempty"`
	Database    string            `json:"database,omitempty"`
	Username    string            `json:"username,omitempty"`
	Password    string            `json:"password,omitempty"`
}

// RDBInstanceInfo holds info about a running database instance
type RDBInstanceInfo struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Type        RDBType           `json:"type"`
	Status      string            `json:"status"`
	Port        string            `json:"port"`
	Host        string            `json:"host"`
	Database    string            `json:"database"`
	Username    string            `json:"username"`
	Volumes     []RDBVolume       `json:"volumes"`
	Environment map[string]string `json:"environment"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}
