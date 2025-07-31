package models

import (
	"time"

	"github.com/google/uuid"
)

// ResourceType defines the type of resource being monitored
type ResourceType string

const (
	ResourceTypeCompute ResourceType = "compute"
	ResourceTypeLambda  ResourceType = "lambda"
	ResourceTypeStorage ResourceType = "storage"
	ResourceTypeRDB     ResourceType = "rdb"
	ResourceTypeSecret  ResourceType = "secret"
)

// ResourceStatus defines the current status of a resource
type ResourceStatus string

const (
	ResourceStatusActive   ResourceStatus = "active"
	ResourceStatusInactive ResourceStatus = "inactive"
	ResourceStatusDeleted  ResourceStatus = "deleted"
)

// ResourceUsage represents a single usage record for a resource
// swagger:model
type ResourceUsage struct {
	// The unique ID of the usage record
	ID uuid.UUID `json:"id" gorm:"type:text;primaryKey"`
	// The creation timestamp
	CreatedAt time.Time `json:"created_at"`
	// The last update timestamp
	UpdatedAt time.Time `json:"updated_at"`

	// The ID of the account that owns this resource
	AccountID uuid.UUID `json:"account_id" gorm:"type:text;not null;index"`
	// The ID of the organization this account belongs to
	OrganizationID uuid.UUID `json:"organization_id" gorm:"type:text;not null;index"`
	// The type of resource (compute, lambda, storage, rdb, secret)
	ResourceType ResourceType `json:"resource_type" gorm:"type:text;not null;index"`
	// The unique identifier of the resource
	ResourceID string `json:"resource_id" gorm:"type:text;not null;index"`
	// The name of the resource
	ResourceName string `json:"resource_name" gorm:"type:text;not null"`
	// The current status of the resource
	Status ResourceStatus `json:"status" gorm:"type:text;not null;index"`

	// Usage metrics
	// CPU usage in percentage (0-100)
	CPUUsage float64 `json:"cpu_usage" gorm:"type:real;default:0"`
	// Memory usage in MB
	MemoryUsage float64 `json:"memory_usage" gorm:"type:real;default:0"`
	// Storage usage in MB
	StorageUsage float64 `json:"storage_usage" gorm:"type:real;default:0"`
	// Network usage in MB
	NetworkUsage float64 `json:"network_usage" gorm:"type:real;default:0"`
	// Number of requests/executions
	RequestCount int64 `json:"request_count" gorm:"type:integer;default:0"`
	// Execution time in seconds
	ExecutionTime float64 `json:"execution_time" gorm:"type:real;default:0"`

	// Billing information
	// Cost per hour in USD
	CostPerHour float64 `json:"cost_per_hour" gorm:"type:real;default:0"`
	// Total cost for this usage period
	TotalCost float64 `json:"total_cost" gorm:"type:real;default:0"`
	// Currency for billing (default: USD)
	Currency string `json:"currency" gorm:"type:text;default:'USD'"`

	// Timestamps for usage tracking
	// When the resource was created
	ResourceCreatedAt time.Time `json:"resource_created_at"`
	// When the resource was last active
	LastActiveAt *time.Time `json:"last_active_at"`
	// When this usage record was last updated
	LastUpdatedAt time.Time `json:"last_updated_at"`

	// Additional metadata
	// Resource configuration (JSON)
	Configuration string `json:"configuration" gorm:"type:text"`
	// Resource tags (JSON)
	Tags string `json:"tags" gorm:"type:text"`
	// Notes or comments
	Notes string `json:"notes" gorm:"type:text"`

	// GORM Associations
	Account      Account      `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
}

// ResourceUsageHistory represents historical usage data for billing and analytics
// swagger:model
type ResourceUsageHistory struct {
	// The unique ID of the history record
	ID uuid.UUID `json:"id" gorm:"type:text;primaryKey"`
	// The creation timestamp
	CreatedAt time.Time `json:"created_at"`

	// The ID of the account that owns this resource
	AccountID uuid.UUID `json:"account_id" gorm:"type:text;not null;index"`
	// The ID of the organization this account belongs to
	OrganizationID uuid.UUID `json:"organization_id" gorm:"type:text;not null;index"`
	// The type of resource
	ResourceType ResourceType `json:"resource_type" gorm:"type:text;not null;index"`
	// The unique identifier of the resource
	ResourceID string `json:"resource_id" gorm:"type:text;not null;index"`
	// The name of the resource
	ResourceName string `json:"resource_name" gorm:"type:text;not null"`

	// Historical metrics
	// Average CPU usage for the period
	AvgCPUUsage float64 `json:"avg_cpu_usage" gorm:"type:real"`
	// Peak CPU usage for the period
	PeakCPUUsage float64 `json:"peak_cpu_usage" gorm:"type:real"`
	// Average memory usage for the period
	AvgMemoryUsage float64 `json:"avg_memory_usage" gorm:"type:real"`
	// Peak memory usage for the period
	PeakMemoryUsage float64 `json:"peak_memory_usage" gorm:"type:real"`
	// Total storage usage for the period
	TotalStorageUsage float64 `json:"total_storage_usage" gorm:"type:real"`
	// Total network usage for the period
	TotalNetworkUsage float64 `json:"total_network_usage" gorm:"type:real"`
	// Total requests for the period
	TotalRequests int64 `json:"total_requests" gorm:"type:integer"`
	// Total execution time for the period
	TotalExecutionTime float64 `json:"total_execution_time" gorm:"type:real"`

	// Billing information
	// Total cost for the period
	TotalCost float64 `json:"total_cost" gorm:"type:real"`
	// Currency for billing
	Currency string `json:"currency" gorm:"type:text;default:'USD'"`

	// Time period
	// Start of the period
	PeriodStart time.Time `json:"period_start" gorm:"type:datetime;not null;index"`
	// End of the period
	PeriodEnd time.Time `json:"period_end" gorm:"type:datetime;not null;index"`
	// Duration in hours
	DurationHours float64 `json:"duration_hours" gorm:"type:real"`

	// GORM Associations
	Account      Account      `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
}

// BillingPeriod represents a billing period for an organization
// swagger:model
type BillingPeriod struct {
	// The unique ID of the billing period
	ID uuid.UUID `json:"id" gorm:"type:text;primaryKey"`
	// The creation timestamp
	CreatedAt time.Time `json:"created_at"`
	// The last update timestamp
	UpdatedAt time.Time `json:"updated_at"`

	// The ID of the organization
	OrganizationID uuid.UUID `json:"organization_id" gorm:"type:text;not null;index"`
	// The start of the billing period
	PeriodStart time.Time `json:"period_start" gorm:"type:datetime;not null;index"`
	// The end of the billing period
	PeriodEnd time.Time `json:"period_end" gorm:"type:datetime;not null;index"`
	// The total cost for the period
	TotalCost float64 `json:"total_cost" gorm:"type:real;default:0"`
	// The currency for billing
	Currency string `json:"currency" gorm:"type:text;default:'USD'"`
	// Whether the bill has been paid
	IsPaid bool `json:"is_paid" gorm:"type:boolean;default:false"`
	// The date when the bill was paid
	PaidAt *time.Time `json:"paid_at"`
	// The payment method used
	PaymentMethod string `json:"payment_method" gorm:"type:text"`
	// The transaction ID for the payment
	TransactionID string `json:"transaction_id" gorm:"type:text"`

	// Breakdown by resource type
	ComputeCost float64 `json:"compute_cost" gorm:"type:real;default:0"`
	LambdaCost  float64 `json:"lambda_cost" gorm:"type:real;default:0"`
	StorageCost float64 `json:"storage_cost" gorm:"type:real;default:0"`
	RDBCost     float64 `json:"rdb_cost" gorm:"type:real;default:0"`
	SecretCost  float64 `json:"secret_cost" gorm:"type:real;default:0"`

	// GORM Associations
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
}

// MonitoringMetrics represents real-time monitoring metrics
// swagger:model
type MonitoringMetrics struct {
	// The unique ID of the metrics record
	ID uuid.UUID `json:"id" gorm:"type:text;primaryKey"`
	// The creation timestamp
	CreatedAt time.Time `json:"created_at"`

	// The ID of the organization
	OrganizationID uuid.UUID `json:"organization_id" gorm:"type:text;not null;index"`
	// The type of resource
	ResourceType ResourceType `json:"resource_type" gorm:"type:text;not null;index"`
	// The unique identifier of the resource
	ResourceID string `json:"resource_id" gorm:"type:text;not null;index"`

	// Real-time metrics
	// Current CPU usage
	CPUUsage float64 `json:"cpu_usage" gorm:"type:real"`
	// Current memory usage
	MemoryUsage float64 `json:"memory_usage" gorm:"type:real"`
	// Current storage usage
	StorageUsage float64 `json:"storage_usage" gorm:"type:real"`
	// Current network usage
	NetworkUsage float64 `json:"network_usage" gorm:"type:real"`
	// Current request rate (requests per second)
	RequestRate float64 `json:"request_rate" gorm:"type:real"`
	// Current error rate (errors per second)
	ErrorRate float64 `json:"error_rate" gorm:"type:real"`
	// Current response time in milliseconds
	ResponseTime float64 `json:"response_time" gorm:"type:real"`

	// Status information
	// Current status of the resource
	Status ResourceStatus `json:"status" gorm:"type:text"`
	// Health check status
	HealthStatus string `json:"health_status" gorm:"type:text"`
	// Last health check timestamp
	LastHealthCheck *time.Time `json:"last_health_check"`

	// GORM Associations
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
}

// ResourceUsageRequest represents a request to create or update resource usage
type ResourceUsageRequest struct {
	ResourceType ResourceType   `json:"resource_type" binding:"required"`
	ResourceID   string         `json:"resource_id" binding:"required"`
	ResourceName string         `json:"resource_name" binding:"required"`
	Status       ResourceStatus `json:"status" binding:"required"`

	// Optional metrics
	CPUUsage      *float64 `json:"cpu_usage,omitempty"`
	MemoryUsage   *float64 `json:"memory_usage,omitempty"`
	StorageUsage  *float64 `json:"storage_usage,omitempty"`
	NetworkUsage  *float64 `json:"network_usage,omitempty"`
	RequestCount  *int64   `json:"request_count,omitempty"`
	ExecutionTime *float64 `json:"execution_time,omitempty"`

	// Optional billing
	CostPerHour *float64 `json:"cost_per_hour,omitempty"`
	Currency    *string  `json:"currency,omitempty"`

	// Optional metadata
	Configuration *string `json:"configuration,omitempty"`
	Tags          *string `json:"tags,omitempty"`
	Notes         *string `json:"notes,omitempty"`
}

// ResourceUsageResponse represents the response for resource usage queries
type ResourceUsageResponse struct {
	CurrentUsage    *ResourceUsage         `json:"current_usage,omitempty"`
	HistoricalUsage []ResourceUsageHistory `json:"historical_usage,omitempty"`
	Metrics         *MonitoringMetrics     `json:"metrics,omitempty"`
	BillingPeriod   *BillingPeriod         `json:"billing_period,omitempty"`
	Summary         ResourceUsageSummary   `json:"summary"`
}

// ResourceUsageSummary represents a summary of resource usage
type ResourceUsageSummary struct {
	TotalResources  int                      `json:"total_resources"`
	ActiveResources int                      `json:"active_resources"`
	TotalCost       float64                  `json:"total_cost"`
	Currency        string                   `json:"currency"`
	UsageByType     map[ResourceType]int     `json:"usage_by_type"`
	CostByType      map[ResourceType]float64 `json:"cost_by_type"`
	MonthlyTrend    []MonthlyUsageTrend      `json:"monthly_trend,omitempty"`
}

// MonthlyUsageTrend represents monthly usage trends
type MonthlyUsageTrend struct {
	Month     string  `json:"month"`
	TotalCost float64 `json:"total_cost"`
	Resources int     `json:"resources"`
}
