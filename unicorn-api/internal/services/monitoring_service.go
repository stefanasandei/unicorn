package services

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"

	"unicorn-api/internal/models"
	"unicorn-api/internal/stores"
)

// MonitoringService provides monitoring and billing functionality
type MonitoringService struct {
	monitoringStore stores.MonitoringStore
	iamStore        stores.IAMStore
}

// NewMonitoringService creates a new monitoring service
func NewMonitoringService(monitoringStore stores.MonitoringStore, iamStore stores.IAMStore) *MonitoringService {
	return &MonitoringService{
		monitoringStore: monitoringStore,
		iamStore:        iamStore,
	}
}

// Pricing constants (in USD per hour)
const (
	ComputeMicroCostPerHour = 0.05 // $0.05/hour for micro compute
	ComputeSmallCostPerHour = 0.10 // $0.10/hour for small compute
	LambdaCostPerHour       = 0.02 // $0.02/hour for lambda executions
	StorageCostPerGBPerHour = 0.01 // $0.01/GB/hour for storage
	RDBCostPerHour          = 0.15 // $0.15/hour for RDB instances
	SecretCostPerHour       = 0.01 // $0.01/hour for secrets
)

// TrackResourceCreation tracks when a new resource is created
func (s *MonitoringService) TrackResourceCreation(
	accountID, organizationID uuid.UUID,
	resourceType models.ResourceType,
	resourceID, resourceName string,
	configuration string,
) error {
	// Verify account exists
	_, err := s.iamStore.GetAccountByID(accountID.String())
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}

	// Calculate cost per hour based on resource type
	costPerHour := s.calculateCostPerHour(resourceType, configuration)

	// Create resource usage record
	usage := &models.ResourceUsage{
		AccountID:         accountID,
		OrganizationID:    organizationID,
		ResourceType:      resourceType,
		ResourceID:        resourceID,
		ResourceName:      resourceName,
		Status:            models.ResourceStatusActive,
		CostPerHour:       costPerHour,
		Currency:          "USD",
		ResourceCreatedAt: time.Now(),
		LastActiveAt:      &time.Time{},
		LastUpdatedAt:     time.Now(),
		Configuration:     configuration,
	}

	// Set initial LastActiveAt to now
	now := time.Now()
	usage.LastActiveAt = &now

	return s.monitoringStore.CreateResourceUsage(usage)
}

// TrackResourceUpdate updates resource usage when a resource is modified
func (s *MonitoringService) TrackResourceUpdate(
	resourceID string,
	resourceType models.ResourceType,
	status models.ResourceStatus,
	metrics map[string]interface{},
) error {
	// Get existing usage record
	usage, err := s.monitoringStore.GetResourceUsageByResourceID(resourceID, resourceType)
	if err != nil {
		return fmt.Errorf("failed to get resource usage: %w", err)
	}

	// Update status
	usage.Status = status
	usage.LastUpdatedAt = time.Now()

	// Update metrics if provided
	if metrics != nil {
		if cpuUsage, ok := metrics["cpu_usage"].(float64); ok {
			usage.CPUUsage = cpuUsage
		}
		if memoryUsage, ok := metrics["memory_usage"].(float64); ok {
			usage.MemoryUsage = memoryUsage
		}
		if storageUsage, ok := metrics["storage_usage"].(float64); ok {
			usage.StorageUsage = storageUsage
		}
		if networkUsage, ok := metrics["network_usage"].(float64); ok {
			usage.NetworkUsage = networkUsage
		}
		if requestCount, ok := metrics["request_count"].(int64); ok {
			usage.RequestCount = requestCount
		}
		if executionTime, ok := metrics["execution_time"].(float64); ok {
			usage.ExecutionTime = executionTime
		}
	}

	// Update last active time if resource is active
	if status == models.ResourceStatusActive {
		now := time.Now()
		usage.LastActiveAt = &now
	}

	// Calculate total cost based on usage duration
	usage.TotalCost = s.calculateTotalCost(usage)

	return s.monitoringStore.UpdateResourceUsage(usage)
}

// TrackResourceDeletion marks a resource as deleted
func (s *MonitoringService) TrackResourceDeletion(
	resourceID string,
	resourceType models.ResourceType,
) error {
	// Get existing usage record
	usage, err := s.monitoringStore.GetResourceUsageByResourceID(resourceID, resourceType)
	if err != nil {
		return fmt.Errorf("failed to get resource usage: %w", err)
	}

	// Mark as deleted
	usage.Status = models.ResourceStatusDeleted
	usage.LastUpdatedAt = time.Now()

	// Calculate final cost
	usage.TotalCost = s.calculateTotalCost(usage)

	// Update the record
	err = s.monitoringStore.UpdateResourceUsage(usage)
	if err != nil {
		return fmt.Errorf("failed to update resource usage: %w", err)
	}

	// Create historical record
	history := &models.ResourceUsageHistory{
		AccountID:          usage.AccountID,
		OrganizationID:     usage.OrganizationID,
		ResourceType:       usage.ResourceType,
		ResourceID:         usage.ResourceID,
		ResourceName:       usage.ResourceName,
		AvgCPUUsage:        usage.CPUUsage,
		PeakCPUUsage:       usage.CPUUsage,
		AvgMemoryUsage:     usage.MemoryUsage,
		PeakMemoryUsage:    usage.MemoryUsage,
		TotalStorageUsage:  usage.StorageUsage,
		TotalNetworkUsage:  usage.NetworkUsage,
		TotalRequests:      usage.RequestCount,
		TotalExecutionTime: usage.ExecutionTime,
		TotalCost:          usage.TotalCost,
		Currency:           usage.Currency,
		PeriodStart:        usage.ResourceCreatedAt,
		PeriodEnd:          time.Now(),
		DurationHours:      time.Since(usage.ResourceCreatedAt).Hours(),
	}

	return s.monitoringStore.CreateResourceUsageHistory(history)
}

// GetResourceUsageSummary gets a summary of resource usage for an organization
func (s *MonitoringService) GetResourceUsageSummary(organizationID string, start, end time.Time) (*models.ResourceUsageResponse, error) {
	// Get current usage
	currentUsage, err := s.monitoringStore.GetResourceUsageByOrganization(organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current usage: %w", err)
	}

	// Get historical usage
	historicalUsage, err := s.monitoringStore.GetResourceUsageHistoryByOrganization(organizationID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get historical usage: %w", err)
	}

	// Get summary
	summary, err := s.monitoringStore.GetResourceUsageSummary(organizationID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage summary: %w", err)
	}

	// Get current billing period
	billingPeriod, err := s.monitoringStore.GetCurrentBillingPeriod(organizationID)
	if err != nil {
		// If no current billing period, create one
		billingPeriod = s.createBillingPeriod(organizationID)
	}

	response := &models.ResourceUsageResponse{
		CurrentUsage:    nil, // Will be set below if available
		HistoricalUsage: historicalUsage,
		BillingPeriod:   billingPeriod,
		Summary:         *summary,
	}

	// Set current usage if available
	if len(currentUsage) > 0 {
		response.CurrentUsage = &currentUsage[0]
	}

	return response, nil
}

// GetMonitoringMetrics gets real-time monitoring metrics for a resource
func (s *MonitoringService) GetMonitoringMetrics(resourceID string, resourceType models.ResourceType) (*models.MonitoringMetrics, error) {
	metrics, err := s.monitoringStore.GetMonitoringMetricsByResource(resourceID, resourceType)
	if err != nil {
		return nil, fmt.Errorf("failed to get monitoring metrics: %w", err)
	}
	return metrics, nil
}

// UpdateMonitoringMetrics updates real-time monitoring metrics
func (s *MonitoringService) UpdateMonitoringMetrics(
	organizationID uuid.UUID,
	resourceID string,
	resourceType models.ResourceType,
	metrics map[string]interface{},
) error {
	// Try to get existing metrics
	existingMetrics, err := s.monitoringStore.GetMonitoringMetricsByResource(resourceID, resourceType)
	if err != nil {
		// Create new metrics if not found
		existingMetrics = &models.MonitoringMetrics{
			OrganizationID: organizationID,
			ResourceType:   resourceType,
			ResourceID:     resourceID,
		}
	}

	// Update metrics
	if cpuUsage, ok := metrics["cpu_usage"].(float64); ok {
		existingMetrics.CPUUsage = cpuUsage
	}
	if memoryUsage, ok := metrics["memory_usage"].(float64); ok {
		existingMetrics.MemoryUsage = memoryUsage
	}
	if storageUsage, ok := metrics["storage_usage"].(float64); ok {
		existingMetrics.StorageUsage = storageUsage
	}
	if networkUsage, ok := metrics["network_usage"].(float64); ok {
		existingMetrics.NetworkUsage = networkUsage
	}
	if requestRate, ok := metrics["request_rate"].(float64); ok {
		existingMetrics.RequestRate = requestRate
	}
	if errorRate, ok := metrics["error_rate"].(float64); ok {
		existingMetrics.ErrorRate = errorRate
	}
	if responseTime, ok := metrics["response_time"].(float64); ok {
		existingMetrics.ResponseTime = responseTime
	}
	if status, ok := metrics["status"].(models.ResourceStatus); ok {
		existingMetrics.Status = status
	}
	if healthStatus, ok := metrics["health_status"].(string); ok {
		existingMetrics.HealthStatus = healthStatus
	}

	// Update health check timestamp
	now := time.Now()
	existingMetrics.LastHealthCheck = &now

	if existingMetrics.ID == uuid.Nil {
		return s.monitoringStore.CreateMonitoringMetrics(existingMetrics)
	}
	return s.monitoringStore.UpdateMonitoringMetrics(existingMetrics)
}

// GenerateMonthlyBilling generates monthly billing for an organization
func (s *MonitoringService) GenerateMonthlyBilling(organizationID string, year int, month int) (*models.BillingPeriod, error) {
	// Get monthly usage history
	usageHistory, err := s.monitoringStore.GetMonthlyUsageHistory(organizationID, year, month)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly usage history: %w", err)
	}

	// Calculate total costs by type
	computeCost := 0.0
	lambdaCost := 0.0
	storageCost := 0.0
	rdbCost := 0.0
	secretCost := 0.0

	for _, usage := range usageHistory {
		switch usage.ResourceType {
		case models.ResourceTypeCompute:
			computeCost += usage.TotalCost
		case models.ResourceTypeLambda:
			lambdaCost += usage.TotalCost
		case models.ResourceTypeStorage:
			storageCost += usage.TotalCost
		case models.ResourceTypeRDB:
			rdbCost += usage.TotalCost
		case models.ResourceTypeSecret:
			secretCost += usage.TotalCost
		}
	}

	totalCost := computeCost + lambdaCost + storageCost + rdbCost + secretCost

	// Create billing period
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	orgUUID, err := uuid.Parse(organizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	billingPeriod := &models.BillingPeriod{
		OrganizationID: orgUUID,
		PeriodStart:    startDate,
		PeriodEnd:      endDate,
		TotalCost:      totalCost,
		Currency:       "USD",
		IsPaid:         false,
		ComputeCost:    computeCost,
		LambdaCost:     lambdaCost,
		StorageCost:    storageCost,
		RDBCost:        rdbCost,
		SecretCost:     secretCost,
	}

	err = s.monitoringStore.CreateBillingPeriod(billingPeriod)
	if err != nil {
		return nil, fmt.Errorf("failed to create billing period: %w", err)
	}

	return billingPeriod, nil
}

// GetBillingHistory gets billing history for an organization
func (s *MonitoringService) GetBillingHistory(organizationID string) ([]models.BillingPeriod, error) {
	return s.monitoringStore.GetBillingPeriodsByOrganization(organizationID)
}

// GetActiveResources gets all active resources for an organization
func (s *MonitoringService) GetActiveResources(organizationID string) ([]models.ResourceUsage, error) {
	orgUUID, err := uuid.Parse(organizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	activeResources, err := s.monitoringStore.GetActiveResourcesByOrganization(orgUUID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get active resources: %w", err)
	}

	// Update costs and LastActiveAt for each resource
	now := time.Now()
	for i := range activeResources {
		// Update LastActiveAt for active resources
		activeResources[i].LastActiveAt = &now
		activeResources[i].TotalCost = s.calculateTotalCost(&activeResources[i])
	}

	return activeResources, nil
}

// GetMonthlyUsageTrends gets monthly usage trends for an organization
func (s *MonitoringService) GetMonthlyUsageTrends(organizationID string, months int) ([]models.MonthlyUsageTrend, error) {
	return s.monitoringStore.GetMonthlyUsageTrends(organizationID, months)
}

// calculateCostPerHour calculates the cost per hour for a resource type
func (s *MonitoringService) calculateCostPerHour(resourceType models.ResourceType, configuration string) float64 {
	switch resourceType {
	case models.ResourceTypeCompute:
		// Parse configuration to determine preset
		if configuration != "" {
			// For now, use default pricing. In a real implementation,
			// you would parse the configuration to determine the preset
			return ComputeMicroCostPerHour
		}
		return ComputeMicroCostPerHour
	case models.ResourceTypeLambda:
		return LambdaCostPerHour
	case models.ResourceTypeStorage:
		return StorageCostPerGBPerHour
	case models.ResourceTypeRDB:
		return RDBCostPerHour
	case models.ResourceTypeSecret:
		return SecretCostPerHour
	default:
		return 0.0
	}
}

// calculateTotalCost calculates the total cost for a resource based on usage duration
func (s *MonitoringService) calculateTotalCost(usage *models.ResourceUsage) float64 {
	var endTime time.Time
	if usage.LastActiveAt == nil {
		// If LastActiveAt is nil, use current time for active resources
		endTime = time.Now()
	} else {
		endTime = *usage.LastActiveAt
	}

	duration := endTime.Sub(usage.ResourceCreatedAt)
	hours := duration.Hours()

	// Calculate base cost
	baseCost := usage.CostPerHour * hours

	// Add additional costs based on usage metrics
	storageCost := (usage.StorageUsage / 1024.0) * StorageCostPerGBPerHour * hours // Convert MB to GB
	networkCost := (usage.NetworkUsage / 1024.0) * 0.001 * hours                   // $0.001 per GB

	return math.Round((baseCost+storageCost+networkCost)*100) / 100 // Round to 2 decimal places
}

// createBillingPeriod creates a new billing period for an organization
func (s *MonitoringService) createBillingPeriod(organizationID string) *models.BillingPeriod {
	orgUUID, _ := uuid.Parse(organizationID)
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	return &models.BillingPeriod{
		OrganizationID: orgUUID,
		PeriodStart:    startDate,
		PeriodEnd:      endDate,
		TotalCost:      0.0,
		Currency:       "USD",
		IsPaid:         false,
	}
}
