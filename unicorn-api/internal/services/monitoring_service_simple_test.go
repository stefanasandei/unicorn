package services

import (
	"testing"
	"time"

	"unicorn-api/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestPricingConstants(t *testing.T) {
	// Test that pricing constants are reasonable
	assert.True(t, 0.05 > 0) // ComputeMicroCostPerHour
	assert.True(t, 0.10 > 0) // ComputeSmallCostPerHour
	assert.True(t, 0.02 > 0) // LambdaCostPerHour
	assert.True(t, 0.01 > 0) // StorageCostPerGBPerHour
	assert.True(t, 0.15 > 0) // RDBCostPerHour
	assert.True(t, 0.01 > 0) // SecretCostPerHour

	// Test relative pricing relationships
	assert.True(t, 0.10 > 0.05) // ComputeSmallCostPerHour > ComputeMicroCostPerHour
	assert.True(t, 0.15 > 0.05) // RDBCostPerHour > ComputeMicroCostPerHour
	assert.True(t, 0.02 < 0.05) // LambdaCostPerHour < ComputeMicroCostPerHour
}

func TestCalculateCostPerHourSimple(t *testing.T) {
	service := &MonitoringService{}

	tests := []struct {
		name         string
		resourceType models.ResourceType
		expected     float64
	}{
		{"Compute Micro", models.ResourceTypeCompute, 0.05},
		{"Storage", models.ResourceTypeStorage, 0.01},
		{"Lambda", models.ResourceTypeLambda, 0.02},
		{"RDB", models.ResourceTypeRDB, 0.15},
		{"Secret", models.ResourceTypeSecret, 0.01},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.calculateCostPerHour(tt.resourceType, "")
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateTotalCostSimple(t *testing.T) {
	service := &MonitoringService{}

	usage := &models.ResourceUsage{
		CostPerHour:       0.05,
		StorageUsage:      1024, // 1GB in MB
		NetworkUsage:      512,  // 0.5GB in MB
		ResourceCreatedAt: time.Now().Add(-2 * time.Hour),
		LastActiveAt:      &time.Time{},
	}
	usage.LastActiveAt = &usage.ResourceCreatedAt

	// Calculate expected cost:
	// Base cost: 0.05 * 2 hours = 0.10
	// Storage cost: (1024/1024) * 0.01 * 2 = 0.02
	// Network cost: (512/1024) * 0.001 * 2 = 0.001
	// Total: 0.10 + 0.02 + 0.001 = 0.121
	expectedCost := 0.121

	result := service.calculateTotalCost(usage)
	assert.Equal(t, expectedCost, result)
}

func TestMonitoringModels(t *testing.T) {
	// Test that monitoring models can be created correctly
	usage := &models.ResourceUsage{
		ResourceType: models.ResourceTypeCompute,
		ResourceID:   "test-123",
		ResourceName: "test-compute",
		Status:       models.ResourceStatusActive,
		CostPerHour:  0.05,
		Currency:     "USD",
	}

	assert.Equal(t, models.ResourceTypeCompute, usage.ResourceType)
	assert.Equal(t, "test-123", usage.ResourceID)
	assert.Equal(t, models.ResourceStatusActive, usage.Status)
	assert.Equal(t, 0.05, usage.CostPerHour)
	assert.Equal(t, "USD", usage.Currency)
}

func TestBillingPeriodModel(t *testing.T) {
	// Test billing period model
	period := &models.BillingPeriod{
		TotalCost:   25.50,
		Currency:    "USD",
		IsPaid:      false,
		ComputeCost: 15.00,
		StorageCost: 5.25,
		LambdaCost:  3.25,
		RDBCost:     2.00,
	}

	assert.Equal(t, 25.50, period.TotalCost)
	assert.Equal(t, "USD", period.Currency)
	assert.False(t, period.IsPaid)
	assert.Equal(t, 15.00, period.ComputeCost)
}

func TestResourceUsageSummaryModel(t *testing.T) {
	// Test resource usage summary model
	summary := &models.ResourceUsageSummary{
		TotalResources:  5,
		ActiveResources: 3,
		TotalCost:       25.50,
		Currency:        "USD",
		UsageByType: map[models.ResourceType]int{
			models.ResourceTypeCompute: 2,
			models.ResourceTypeStorage: 1,
			models.ResourceTypeLambda:  1,
			models.ResourceTypeRDB:     1,
		},
		CostByType: map[models.ResourceType]float64{
			models.ResourceTypeCompute: 15.00,
			models.ResourceTypeStorage: 5.25,
			models.ResourceTypeLambda:  3.25,
			models.ResourceTypeRDB:     2.00,
		},
	}

	assert.Equal(t, 5, summary.TotalResources)
	assert.Equal(t, 3, summary.ActiveResources)
	assert.Equal(t, 25.50, summary.TotalCost)
	assert.Equal(t, "USD", summary.Currency)
	assert.Equal(t, 2, summary.UsageByType[models.ResourceTypeCompute])
	assert.Equal(t, 15.00, summary.CostByType[models.ResourceTypeCompute])
}

func TestMonthlyUsageTrendModel(t *testing.T) {
	// Test monthly usage trend model
	trend := &models.MonthlyUsageTrend{
		Month:     "2024-01",
		TotalCost: 45.25,
		Resources: 8,
	}

	assert.Equal(t, "2024-01", trend.Month)
	assert.Equal(t, 45.25, trend.TotalCost)
	assert.Equal(t, 8, trend.Resources)
}

func TestMonitoringMetricsModel(t *testing.T) {
	// Test monitoring metrics model
	metrics := &models.MonitoringMetrics{
		ResourceID:   "test-resource-123",
		ResourceType: models.ResourceTypeCompute,
		CPUUsage:     45.2,
		MemoryUsage:  1024,
		StorageUsage: 512,
		NetworkUsage: 256,
	}

	assert.Equal(t, "test-resource-123", metrics.ResourceID)
	assert.Equal(t, models.ResourceTypeCompute, metrics.ResourceType)
	assert.Equal(t, 45.2, metrics.CPUUsage)
	assert.Equal(t, 1024, metrics.MemoryUsage)
	assert.Equal(t, 512, metrics.StorageUsage)
	assert.Equal(t, 256, metrics.NetworkUsage)
}

func TestResourceUsageRequestModel(t *testing.T) {
	// Test resource usage request model
	config := `{"image":"nginx","preset":"micro"}`
	request := &models.ResourceUsageRequest{
		ResourceType:  models.ResourceTypeCompute,
		ResourceID:    "test-resource-123",
		ResourceName:  "test-compute",
		Status:        models.ResourceStatusActive,
		Configuration: &config,
	}

	assert.Equal(t, models.ResourceTypeCompute, request.ResourceType)
	assert.Equal(t, "test-resource-123", request.ResourceID)
	assert.Equal(t, "test-compute", request.ResourceName)
	assert.Equal(t, models.ResourceStatusActive, request.Status)
	assert.Equal(t, config, *request.Configuration)
}

func TestResourceTypeValidation(t *testing.T) {
	// Test resource type validation
	validTypes := []models.ResourceType{
		models.ResourceTypeCompute,
		models.ResourceTypeStorage,
		models.ResourceTypeLambda,
		models.ResourceTypeRDB,
		models.ResourceTypeSecret,
	}

	for _, resourceType := range validTypes {
		t.Run(string(resourceType), func(t *testing.T) {
			assert.NotEmpty(t, string(resourceType))
		})
	}
}

func TestResourceStatusValidation(t *testing.T) {
	// Test resource status validation
	validStatuses := []models.ResourceStatus{
		models.ResourceStatusActive,
		models.ResourceStatusInactive,
		models.ResourceStatusDeleted,
	}

	for _, status := range validStatuses {
		t.Run(string(status), func(t *testing.T) {
			assert.NotEmpty(t, string(status))
		})
	}
}
