package stores

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"unicorn-api/internal/models"
)

// MonitoringStore defines the interface for monitoring data storage
type MonitoringStore interface {
	// Resource Usage Management
	CreateResourceUsage(usage *models.ResourceUsage) error
	UpdateResourceUsage(usage *models.ResourceUsage) error
	GetResourceUsageByID(id string) (*models.ResourceUsage, error)
	GetResourceUsageByResourceID(resourceID string, resourceType models.ResourceType) (*models.ResourceUsage, error)
	GetResourceUsageByAccount(accountID string) ([]models.ResourceUsage, error)
	GetResourceUsageByOrganization(orgID string) ([]models.ResourceUsage, error)
	GetActiveResourcesByOrganization(orgID string) ([]models.ResourceUsage, error)
	DeleteResourceUsage(id string) error

	// Resource Usage History
	CreateResourceUsageHistory(history *models.ResourceUsageHistory) error
	GetResourceUsageHistoryByResource(resourceID string, resourceType models.ResourceType, start, end time.Time) ([]models.ResourceUsageHistory, error)
	GetResourceUsageHistoryByOrganization(orgID string, start, end time.Time) ([]models.ResourceUsageHistory, error)
	GetMonthlyUsageHistory(orgID string, year int, month int) ([]models.ResourceUsageHistory, error)

	// Billing Period Management
	CreateBillingPeriod(period *models.BillingPeriod) error
	UpdateBillingPeriod(period *models.BillingPeriod) error
	GetBillingPeriodByID(id string) (*models.BillingPeriod, error)
	GetBillingPeriodsByOrganization(orgID string) ([]models.BillingPeriod, error)
	GetCurrentBillingPeriod(orgID string) (*models.BillingPeriod, error)
	GetBillingPeriodByDateRange(orgID string, start, end time.Time) (*models.BillingPeriod, error)

	// Monitoring Metrics
	CreateMonitoringMetrics(metrics *models.MonitoringMetrics) error
	UpdateMonitoringMetrics(metrics *models.MonitoringMetrics) error
	GetMonitoringMetricsByResource(resourceID string, resourceType models.ResourceType) (*models.MonitoringMetrics, error)
	GetMonitoringMetricsByOrganization(orgID string) ([]models.MonitoringMetrics, error)
	DeleteMonitoringMetrics(id string) error

	// Analytics and Reporting
	GetResourceUsageSummary(orgID string, start, end time.Time) (*models.ResourceUsageSummary, error)
	GetMonthlyUsageTrends(orgID string, months int) ([]models.MonthlyUsageTrend, error)
	GetCostBreakdownByType(orgID string, start, end time.Time) (map[models.ResourceType]float64, error)
	GetActiveResourceCount(orgID string) (int, error)
	GetTotalCostForPeriod(orgID string, start, end time.Time) (float64, error)

	// Database operations
	DB() *gorm.DB
}

// GORMMonitoringStore implements MonitoringStore using GORM for SQLite
type GORMMonitoringStore struct {
	db *gorm.DB
}

// NewGORMMonitoringStore creates a new GORMMonitoringStore
func NewGORMMonitoringStore(dataSourceName string) (*GORMMonitoringStore, error) {
	db, err := gorm.Open(sqlite.Open(dataSourceName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database with GORM: %w", err)
	}

	// AutoMigrate will create and update tables based on your models
	err = db.AutoMigrate(
		&models.ResourceUsage{},
		&models.ResourceUsageHistory{},
		&models.BillingPeriod{},
		&models.MonitoringMetrics{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate database schema: %w", err)
	}

	return &GORMMonitoringStore{db: db}, nil
}

// CreateResourceUsage inserts a new resource usage record
func (s *GORMMonitoringStore) CreateResourceUsage(usage *models.ResourceUsage) error {
	if usage.ID == uuid.Nil {
		usage.ID = uuid.New()
	}
	usage.CreatedAt = time.Now()
	usage.UpdatedAt = time.Now()
	usage.LastUpdatedAt = time.Now()

	result := s.db.Create(usage)
	if result.Error != nil {
		return fmt.Errorf("failed to create resource usage: %w", result.Error)
	}
	return nil
}

// UpdateResourceUsage updates an existing resource usage record
func (s *GORMMonitoringStore) UpdateResourceUsage(usage *models.ResourceUsage) error {
	usage.UpdatedAt = time.Now()
	usage.LastUpdatedAt = time.Now()

	result := s.db.Save(usage)
	if result.Error != nil {
		return fmt.Errorf("failed to update resource usage: %w", result.Error)
	}
	return nil
}

// GetResourceUsageByID retrieves a resource usage record by ID
func (s *GORMMonitoringStore) GetResourceUsageByID(id string) (*models.ResourceUsage, error) {
	usageID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid usage ID format: %w", err)
	}

	var usage models.ResourceUsage
	result := s.db.Where("id = ?", usageID).First(&usage)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("resource usage not found")
		}
		return nil, fmt.Errorf("failed to get resource usage: %w", result.Error)
	}
	return &usage, nil
}

// GetResourceUsageByResourceID retrieves a resource usage record by resource ID and type
func (s *GORMMonitoringStore) GetResourceUsageByResourceID(resourceID string, resourceType models.ResourceType) (*models.ResourceUsage, error) {
	var usage models.ResourceUsage
	result := s.db.Where("resource_id = ? AND resource_type = ?", resourceID, resourceType).First(&usage)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("resource usage not found")
		}
		return nil, fmt.Errorf("failed to get resource usage: %w", result.Error)
	}
	return &usage, nil
}

// GetResourceUsageByAccount retrieves all resource usage records for an account
func (s *GORMMonitoringStore) GetResourceUsageByAccount(accountID string) ([]models.ResourceUsage, error) {
	accUUID, err := uuid.Parse(accountID)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID format: %w", err)
	}

	var usages []models.ResourceUsage
	result := s.db.Where("account_id = ?", accUUID).Find(&usages)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get resource usage by account: %w", result.Error)
	}
	return usages, nil
}

// GetResourceUsageByOrganization retrieves all resource usage records for an organization
func (s *GORMMonitoringStore) GetResourceUsageByOrganization(orgID string) ([]models.ResourceUsage, error) {
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID format: %w", err)
	}

	var usages []models.ResourceUsage
	result := s.db.Where("organization_id = ?", orgUUID).Find(&usages)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get resource usage by organization: %w", result.Error)
	}
	return usages, nil
}

// GetActiveResourcesByOrganization retrieves all active resource usage records for an organization
func (s *GORMMonitoringStore) GetActiveResourcesByOrganization(orgID string) ([]models.ResourceUsage, error) {
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID format: %w", err)
	}

	var usages []models.ResourceUsage
	result := s.db.Where("organization_id = ? AND status = ?", orgUUID, models.ResourceStatusActive).Find(&usages)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get active resources by organization: %w", result.Error)
	}
	return usages, nil
}

// DeleteResourceUsage deletes a resource usage record
func (s *GORMMonitoringStore) DeleteResourceUsage(id string) error {
	usageID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid usage ID format: %w", err)
	}

	result := s.db.Delete(&models.ResourceUsage{}, usageID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete resource usage: %w", result.Error)
	}
	return nil
}

// CreateResourceUsageHistory inserts a new resource usage history record
func (s *GORMMonitoringStore) CreateResourceUsageHistory(history *models.ResourceUsageHistory) error {
	if history.ID == uuid.Nil {
		history.ID = uuid.New()
	}
	history.CreatedAt = time.Now()

	result := s.db.Create(history)
	if result.Error != nil {
		return fmt.Errorf("failed to create resource usage history: %w", result.Error)
	}
	return nil
}

// GetResourceUsageHistoryByResource retrieves usage history for a specific resource
func (s *GORMMonitoringStore) GetResourceUsageHistoryByResource(resourceID string, resourceType models.ResourceType, start, end time.Time) ([]models.ResourceUsageHistory, error) {
	var history []models.ResourceUsageHistory
	result := s.db.Where("resource_id = ? AND resource_type = ? AND period_start >= ? AND period_end <= ?",
		resourceID, resourceType, start, end).Find(&history)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get resource usage history: %w", result.Error)
	}
	return history, nil
}

// GetResourceUsageHistoryByOrganization retrieves usage history for an organization
func (s *GORMMonitoringStore) GetResourceUsageHistoryByOrganization(orgID string, start, end time.Time) ([]models.ResourceUsageHistory, error) {
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID format: %w", err)
	}

	var history []models.ResourceUsageHistory
	result := s.db.Where("organization_id = ? AND period_start >= ? AND period_end <= ?",
		orgUUID, start, end).Find(&history)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get organization usage history: %w", result.Error)
	}
	return history, nil
}

// GetMonthlyUsageHistory retrieves monthly usage history for an organization
func (s *GORMMonitoringStore) GetMonthlyUsageHistory(orgID string, year int, month int) ([]models.ResourceUsageHistory, error) {
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID format: %w", err)
	}

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	var history []models.ResourceUsageHistory
	result := s.db.Where("organization_id = ? AND period_start >= ? AND period_end <= ?",
		orgUUID, startDate, endDate).Find(&history)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get monthly usage history: %w", result.Error)
	}
	return history, nil
}

// CreateBillingPeriod inserts a new billing period record
func (s *GORMMonitoringStore) CreateBillingPeriod(period *models.BillingPeriod) error {
	if period.ID == uuid.Nil {
		period.ID = uuid.New()
	}
	period.CreatedAt = time.Now()
	period.UpdatedAt = time.Now()

	result := s.db.Create(period)
	if result.Error != nil {
		return fmt.Errorf("failed to create billing period: %w", result.Error)
	}
	return nil
}

// UpdateBillingPeriod updates an existing billing period record
func (s *GORMMonitoringStore) UpdateBillingPeriod(period *models.BillingPeriod) error {
	period.UpdatedAt = time.Now()

	result := s.db.Save(period)
	if result.Error != nil {
		return fmt.Errorf("failed to update billing period: %w", result.Error)
	}
	return nil
}

// GetBillingPeriodByID retrieves a billing period by ID
func (s *GORMMonitoringStore) GetBillingPeriodByID(id string) (*models.BillingPeriod, error) {
	periodID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid billing period ID format: %w", err)
	}

	var period models.BillingPeriod
	result := s.db.Where("id = ?", periodID).First(&period)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("billing period not found")
		}
		return nil, fmt.Errorf("failed to get billing period: %w", result.Error)
	}
	return &period, nil
}

// GetBillingPeriodsByOrganization retrieves all billing periods for an organization
func (s *GORMMonitoringStore) GetBillingPeriodsByOrganization(orgID string) ([]models.BillingPeriod, error) {
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID format: %w", err)
	}

	var periods []models.BillingPeriod
	result := s.db.Where("organization_id = ?", orgUUID).Order("period_start DESC").Find(&periods)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get billing periods by organization: %w", result.Error)
	}
	return periods, nil
}

// GetCurrentBillingPeriod retrieves the current billing period for an organization
func (s *GORMMonitoringStore) GetCurrentBillingPeriod(orgID string) (*models.BillingPeriod, error) {
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID format: %w", err)
	}

	now := time.Now()
	var period models.BillingPeriod
	result := s.db.Where("organization_id = ? AND period_start <= ? AND period_end >= ?",
		orgUUID, now, now).First(&period)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("current billing period not found")
		}
		return nil, fmt.Errorf("failed to get current billing period: %w", result.Error)
	}
	return &period, nil
}

// GetBillingPeriodByDateRange retrieves a billing period for a specific date range
func (s *GORMMonitoringStore) GetBillingPeriodByDateRange(orgID string, start, end time.Time) (*models.BillingPeriod, error) {
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID format: %w", err)
	}

	var period models.BillingPeriod
	result := s.db.Where("organization_id = ? AND period_start = ? AND period_end = ?",
		orgUUID, start, end).First(&period)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("billing period not found for date range")
		}
		return nil, fmt.Errorf("failed to get billing period by date range: %w", result.Error)
	}
	return &period, nil
}

// CreateMonitoringMetrics inserts a new monitoring metrics record
func (s *GORMMonitoringStore) CreateMonitoringMetrics(metrics *models.MonitoringMetrics) error {
	if metrics.ID == uuid.Nil {
		metrics.ID = uuid.New()
	}
	metrics.CreatedAt = time.Now()

	result := s.db.Create(metrics)
	if result.Error != nil {
		return fmt.Errorf("failed to create monitoring metrics: %w", result.Error)
	}
	return nil
}

// UpdateMonitoringMetrics updates an existing monitoring metrics record
func (s *GORMMonitoringStore) UpdateMonitoringMetrics(metrics *models.MonitoringMetrics) error {
	result := s.db.Save(metrics)
	if result.Error != nil {
		return fmt.Errorf("failed to update monitoring metrics: %w", result.Error)
	}
	return nil
}

// GetMonitoringMetricsByResource retrieves monitoring metrics for a specific resource
func (s *GORMMonitoringStore) GetMonitoringMetricsByResource(resourceID string, resourceType models.ResourceType) (*models.MonitoringMetrics, error) {
	var metrics models.MonitoringMetrics
	result := s.db.Where("resource_id = ? AND resource_type = ?", resourceID, resourceType).First(&metrics)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("monitoring metrics not found")
		}
		return nil, fmt.Errorf("failed to get monitoring metrics: %w", result.Error)
	}
	return &metrics, nil
}

// GetMonitoringMetricsByOrganization retrieves all monitoring metrics for an organization
func (s *GORMMonitoringStore) GetMonitoringMetricsByOrganization(orgID string) ([]models.MonitoringMetrics, error) {
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID format: %w", err)
	}

	var metrics []models.MonitoringMetrics
	result := s.db.Where("organization_id = ?", orgUUID).Find(&metrics)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get monitoring metrics by organization: %w", result.Error)
	}
	return metrics, nil
}

// DeleteMonitoringMetrics deletes a monitoring metrics record
func (s *GORMMonitoringStore) DeleteMonitoringMetrics(id string) error {
	metricsID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid metrics ID format: %w", err)
	}

	result := s.db.Delete(&models.MonitoringMetrics{}, metricsID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete monitoring metrics: %w", result.Error)
	}
	return nil
}

// GetResourceUsageSummary retrieves a summary of resource usage for an organization
func (s *GORMMonitoringStore) GetResourceUsageSummary(orgID string, start, end time.Time) (*models.ResourceUsageSummary, error) {
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID format: %w", err)
	}

	// Get current usage
	var currentUsage []models.ResourceUsage
	result := s.db.Where("organization_id = ?", orgUUID).Find(&currentUsage)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get current usage: %w", result.Error)
	}

	// Get historical usage
	var historicalUsage []models.ResourceUsageHistory
	result = s.db.Where("organization_id = ? AND period_start >= ? AND period_end <= ?",
		orgUUID, start, end).Find(&historicalUsage)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get historical usage: %w", result.Error)
	}

	// Calculate summary
	summary := &models.ResourceUsageSummary{
		TotalResources:  len(currentUsage),
		ActiveResources: 0,
		TotalCost:       0,
		Currency:        "USD",
		UsageByType:     make(map[models.ResourceType]int),
		CostByType:      make(map[models.ResourceType]float64),
	}

	// Count active resources and calculate costs
	for _, usage := range currentUsage {
		if usage.Status == models.ResourceStatusActive {
			summary.ActiveResources++
		}
		summary.TotalCost += usage.TotalCost
		summary.UsageByType[usage.ResourceType]++
		summary.CostByType[usage.ResourceType] += usage.TotalCost
	}

	return summary, nil
}

// GetMonthlyUsageTrends retrieves monthly usage trends for an organization
func (s *GORMMonitoringStore) GetMonthlyUsageTrends(orgID string, months int) ([]models.MonthlyUsageTrend, error) {
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID format: %w", err)
	}

	var trends []models.MonthlyUsageTrend
	now := time.Now()

	for i := 0; i < months; i++ {
		month := now.AddDate(0, -i, 0)
		startDate := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

		var totalCost float64
		var resourceCount int64

		// Get total cost for the month
		result := s.db.Model(&models.ResourceUsageHistory{}).
			Where("organization_id = ? AND period_start >= ? AND period_end <= ?",
				orgUUID, startDate, endDate).
			Select("COALESCE(SUM(total_cost), 0) as total_cost, COUNT(*) as resource_count").
			Scan(&struct {
				TotalCost     float64 `json:"total_cost"`
				ResourceCount int64   `json:"resource_count"`
			}{TotalCost: totalCost, ResourceCount: resourceCount})

		if result.Error != nil {
			return nil, fmt.Errorf("failed to get monthly trend: %w", result.Error)
		}

		trends = append(trends, models.MonthlyUsageTrend{
			Month:     month.Format("2006-01"),
			TotalCost: totalCost,
			Resources: int(resourceCount),
		})
	}

	return trends, nil
}

// GetCostBreakdownByType retrieves cost breakdown by resource type
func (s *GORMMonitoringStore) GetCostBreakdownByType(orgID string, start, end time.Time) (map[models.ResourceType]float64, error) {
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID format: %w", err)
	}

	var results []struct {
		ResourceType models.ResourceType `json:"resource_type"`
		TotalCost    float64             `json:"total_cost"`
	}

	result := s.db.Model(&models.ResourceUsageHistory{}).
		Where("organization_id = ? AND period_start >= ? AND period_end <= ?",
			orgUUID, start, end).
		Select("resource_type, SUM(total_cost) as total_cost").
		Group("resource_type").
		Find(&results)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get cost breakdown: %w", result.Error)
	}

	breakdown := make(map[models.ResourceType]float64)
	for _, result := range results {
		breakdown[result.ResourceType] = result.TotalCost
	}

	return breakdown, nil
}

// GetActiveResourceCount retrieves the count of active resources for an organization
func (s *GORMMonitoringStore) GetActiveResourceCount(orgID string) (int, error) {
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		return 0, fmt.Errorf("invalid organization ID format: %w", err)
	}

	var count int64
	result := s.db.Model(&models.ResourceUsage{}).
		Where("organization_id = ? AND status = ?", orgUUID, models.ResourceStatusActive).
		Count(&count)

	if result.Error != nil {
		return 0, fmt.Errorf("failed to get active resource count: %w", result.Error)
	}

	return int(count), nil
}

// GetTotalCostForPeriod retrieves the total cost for a specific period
func (s *GORMMonitoringStore) GetTotalCostForPeriod(orgID string, start, end time.Time) (float64, error) {
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		return 0, fmt.Errorf("invalid organization ID format: %w", err)
	}

	var totalCost float64
	result := s.db.Model(&models.ResourceUsageHistory{}).
		Where("organization_id = ? AND period_start >= ? AND period_end <= ?",
			orgUUID, start, end).
		Select("COALESCE(SUM(total_cost), 0) as total_cost").
		Scan(&totalCost)

	if result.Error != nil {
		return 0, fmt.Errorf("failed to get total cost for period: %w", result.Error)
	}

	return totalCost, nil
}

// DB returns the underlying database connection
func (s *GORMMonitoringStore) DB() *gorm.DB {
	return s.db
}
