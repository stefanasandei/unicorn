package services

import (
	"testing"
	"time"

	"unicorn-api/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockMonitoringStore is a mock implementation of MonitoringStore
type MockMonitoringStore struct {
	mock.Mock
}

func (m *MockMonitoringStore) CreateResourceUsage(usage *models.ResourceUsage) error {
	args := m.Called(usage)
	return args.Error(0)
}

func (m *MockMonitoringStore) UpdateResourceUsage(usage *models.ResourceUsage) error {
	args := m.Called(usage)
	return args.Error(0)
}

func (m *MockMonitoringStore) GetResourceUsageByID(id string) (*models.ResourceUsage, error) {
	args := m.Called(id)
	return args.Get(0).(*models.ResourceUsage), args.Error(1)
}

func (m *MockMonitoringStore) GetResourceUsageByResourceID(resourceID string, resourceType models.ResourceType) (*models.ResourceUsage, error) {
	args := m.Called(resourceID, resourceType)
	return args.Get(0).(*models.ResourceUsage), args.Error(1)
}

func (m *MockMonitoringStore) GetResourceUsageByAccount(accountID string) ([]models.ResourceUsage, error) {
	args := m.Called(accountID)
	return args.Get(0).([]models.ResourceUsage), args.Error(1)
}

func (m *MockMonitoringStore) GetResourceUsageByOrganization(orgID string) ([]models.ResourceUsage, error) {
	args := m.Called(orgID)
	return args.Get(0).([]models.ResourceUsage), args.Error(1)
}

func (m *MockMonitoringStore) GetActiveResourcesByOrganization(orgID string) ([]models.ResourceUsage, error) {
	args := m.Called(orgID)
	return args.Get(0).([]models.ResourceUsage), args.Error(1)
}

func (m *MockMonitoringStore) DeleteResourceUsage(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMonitoringStore) CreateResourceUsageHistory(history *models.ResourceUsageHistory) error {
	args := m.Called(history)
	return args.Error(0)
}

func (m *MockMonitoringStore) GetResourceUsageHistoryByResource(resourceID string, resourceType models.ResourceType, start, end time.Time) ([]models.ResourceUsageHistory, error) {
	args := m.Called(resourceID, resourceType, start, end)
	return args.Get(0).([]models.ResourceUsageHistory), args.Error(1)
}

func (m *MockMonitoringStore) GetResourceUsageHistoryByOrganization(orgID string, start, end time.Time) ([]models.ResourceUsageHistory, error) {
	args := m.Called(orgID, start, end)
	return args.Get(0).([]models.ResourceUsageHistory), args.Error(1)
}

func (m *MockMonitoringStore) GetMonthlyUsageHistory(orgID string, year int, month int) ([]models.ResourceUsageHistory, error) {
	args := m.Called(orgID, year, month)
	return args.Get(0).([]models.ResourceUsageHistory), args.Error(1)
}

func (m *MockMonitoringStore) CreateBillingPeriod(period *models.BillingPeriod) error {
	args := m.Called(period)
	return args.Error(0)
}

func (m *MockMonitoringStore) UpdateBillingPeriod(period *models.BillingPeriod) error {
	args := m.Called(period)
	return args.Error(0)
}

func (m *MockMonitoringStore) GetBillingPeriodByID(id string) (*models.BillingPeriod, error) {
	args := m.Called(id)
	return args.Get(0).(*models.BillingPeriod), args.Error(1)
}

func (m *MockMonitoringStore) GetBillingPeriodsByOrganization(orgID string) ([]models.BillingPeriod, error) {
	args := m.Called(orgID)
	return args.Get(0).([]models.BillingPeriod), args.Error(1)
}

func (m *MockMonitoringStore) GetCurrentBillingPeriod(orgID string) (*models.BillingPeriod, error) {
	args := m.Called(orgID)
	return args.Get(0).(*models.BillingPeriod), args.Error(1)
}

func (m *MockMonitoringStore) GetBillingPeriodByDateRange(orgID string, start, end time.Time) (*models.BillingPeriod, error) {
	args := m.Called(orgID, start, end)
	return args.Get(0).(*models.BillingPeriod), args.Error(1)
}

func (m *MockMonitoringStore) CreateMonitoringMetrics(metrics *models.MonitoringMetrics) error {
	args := m.Called(metrics)
	return args.Error(0)
}

func (m *MockMonitoringStore) UpdateMonitoringMetrics(metrics *models.MonitoringMetrics) error {
	args := m.Called(metrics)
	return args.Error(0)
}

func (m *MockMonitoringStore) GetMonitoringMetricsByResource(resourceID string, resourceType models.ResourceType) (*models.MonitoringMetrics, error) {
	args := m.Called(resourceID, resourceType)
	return args.Get(0).(*models.MonitoringMetrics), args.Error(1)
}

func (m *MockMonitoringStore) GetMonitoringMetricsByOrganization(orgID string) ([]models.MonitoringMetrics, error) {
	args := m.Called(orgID)
	return args.Get(0).([]models.MonitoringMetrics), args.Error(1)
}

func (m *MockMonitoringStore) DeleteMonitoringMetrics(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMonitoringStore) GetResourceUsageSummary(orgID string, start, end time.Time) (*models.ResourceUsageSummary, error) {
	args := m.Called(orgID, start, end)
	return args.Get(0).(*models.ResourceUsageSummary), args.Error(1)
}

func (m *MockMonitoringStore) GetMonthlyUsageTrends(orgID string, months int) ([]models.MonthlyUsageTrend, error) {
	args := m.Called(orgID, months)
	return args.Get(0).([]models.MonthlyUsageTrend), args.Error(1)
}

func (m *MockMonitoringStore) GetCostBreakdownByType(orgID string, start, end time.Time) (map[models.ResourceType]float64, error) {
	args := m.Called(orgID, start, end)
	return args.Get(0).(map[models.ResourceType]float64), args.Error(1)
}

func (m *MockMonitoringStore) GetActiveResourceCount(orgID string) (int, error) {
	args := m.Called(orgID)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockMonitoringStore) GetTotalCostForPeriod(orgID string, start, end time.Time) (float64, error) {
	args := m.Called(orgID, start, end)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockMonitoringStore) DB() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

// MockIAMStore is a mock implementation of IAMStore
type MockIAMStore struct {
	mock.Mock
}

func (m *MockIAMStore) CreateRole(role *models.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockIAMStore) AssignRole(accountID, roleID string) error {
	args := m.Called(accountID, roleID)
	return args.Error(0)
}

func (m *MockIAMStore) GetRoleByName(name string) (*models.Role, error) {
	args := m.Called(name)
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockIAMStore) GetRoleByID(roleID string) (*models.Role, error) {
	args := m.Called(roleID)
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockIAMStore) CreateOrganization(org *models.Organization) error {
	args := m.Called(org)
	return args.Error(0)
}

func (m *MockIAMStore) GetOrganizationByName(name string) (*models.Organization, error) {
	args := m.Called(name)
	return args.Get(0).(*models.Organization), args.Error(1)
}

func (m *MockIAMStore) CreateAccount(account *models.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockIAMStore) GetAccountByEmail(email string) (*models.Account, error) {
	args := m.Called(email)
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockIAMStore) UpdateAccount(account *models.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockIAMStore) GetAccountByID(accountID string) (*models.Account, error) {
	args := m.Called(accountID)
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockIAMStore) GetRolesByOrganizationID(orgID string) ([]models.Role, error) {
	args := m.Called(orgID)
	return args.Get(0).([]models.Role), args.Error(1)
}

func (m *MockIAMStore) GetOrganizationByID(orgID string) (*models.Organization, error) {
	args := m.Called(orgID)
	return args.Get(0).(*models.Organization), args.Error(1)
}

func (m *MockIAMStore) GetAccountsByOrganizationID(orgID string) ([]models.Account, error) {
	args := m.Called(orgID)
	return args.Get(0).([]models.Account), args.Error(1)
}

func TestNewMonitoringService(t *testing.T) {
	mockStore := &MockMonitoringStore{}
	mockIAMStore := &MockIAMStore{}

	service := NewMonitoringService(mockStore, mockIAMStore)

	assert.NotNil(t, service)
	assert.Equal(t, mockStore, service.monitoringStore)
	assert.Equal(t, mockIAMStore, service.iamStore)
}

func TestTrackResourceCreation(t *testing.T) {
	mockStore := &MockMonitoringStore{}
	mockIAMStore := &MockIAMStore{}
	service := NewMonitoringService(mockStore, mockIAMStore)

	userID := uuid.New()
	orgID := uuid.New()
	resourceID := "test-resource-123"
	resourceName := "test-compute"
	configuration := `{"image":"nginx","preset":"micro"}`

	// Mock account lookup
	account := &models.Account{
		ID:             userID,
		OrganizationID: orgID,
	}
	mockIAMStore.On("GetAccountByID", userID.String()).Return(account, nil)

	// Mock resource usage creation
	mockStore.On("CreateResourceUsage", mock.AnythingOfType("*models.ResourceUsage")).Return(nil)

	err := service.TrackResourceCreation(userID, orgID, models.ResourceTypeCompute, resourceID, resourceName, configuration)

	assert.NoError(t, err)
	mockStore.AssertExpectations(t)
	mockIAMStore.AssertExpectations(t)
}

func TestTrackResourceDeletion(t *testing.T) {
	mockStore := &MockMonitoringStore{}
	mockIAMStore := &MockIAMStore{}
	service := NewMonitoringService(mockStore, mockIAMStore)

	resourceID := "test-resource-123"

	// Mock existing resource usage
	existingUsage := &models.ResourceUsage{
		ID:                uuid.New(),
		ResourceID:        resourceID,
		ResourceType:      models.ResourceTypeCompute,
		ResourceName:      "test-compute",
		Status:            models.ResourceStatusActive,
		CostPerHour:       0.05,
		Currency:          "USD",
		ResourceCreatedAt: time.Now().Add(-2 * time.Hour),
		LastActiveAt:      &time.Time{},
		LastUpdatedAt:     time.Now(),
	}

	mockStore.On("GetResourceUsageByResourceID", resourceID, models.ResourceTypeCompute).Return(existingUsage, nil)

	// Mock resource usage update
	mockStore.On("UpdateResourceUsage", mock.AnythingOfType("*models.ResourceUsage")).Return(nil)

	// Mock resource usage history creation
	mockStore.On("CreateResourceUsageHistory", mock.AnythingOfType("*models.ResourceUsageHistory")).Return(nil)

	err := service.TrackResourceDeletion(resourceID, models.ResourceTypeCompute)

	assert.NoError(t, err)
	mockStore.AssertExpectations(t)
}

func TestCalculateCostPerHour(t *testing.T) {
	service := &MonitoringService{}

	tests := []struct {
		name          string
		resourceType  models.ResourceType
		configuration string
		expected      float64
	}{
		{"Compute Micro", models.ResourceTypeCompute, `{"preset":"micro"}`, ComputeMicroCostPerHour},
		{"Storage", models.ResourceTypeStorage, "", StorageCostPerGBPerHour},
		{"Lambda", models.ResourceTypeLambda, "", LambdaCostPerHour},
		{"RDB", models.ResourceTypeRDB, "", RDBCostPerHour},
		{"Secret", models.ResourceTypeSecret, "", SecretCostPerHour},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.calculateCostPerHour(tt.resourceType, tt.configuration)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateTotalCost(t *testing.T) {
	service := &MonitoringService{}

	// Use fixed times for deterministic testing
	createdAt := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	lastActive := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC) // 2 hours later

	usage := &models.ResourceUsage{
		CostPerHour:       0.05,
		StorageUsage:      1024, // 1GB in MB
		NetworkUsage:      512,  // 0.5GB in MB
		ResourceCreatedAt: createdAt,
		LastActiveAt:      &lastActive,
	}

	// Calculate expected cost:
	// Base cost: 0.05 * 2 hours = 0.10
	// Storage cost: (1024/1024) * 0.01 * 2 = 0.02
	// Network cost: (512/1024) * 0.001 * 2 = 0.001
	// Total: 0.10 + 0.02 + 0.001 = 0.121
	// But due to rounding in the implementation, it becomes 0.12
	expectedCost := 0.12

	result := service.calculateTotalCost(usage)
	assert.Equal(t, expectedCost, result)
}

func TestGetResourceUsageSummary(t *testing.T) {
	mockStore := &MockMonitoringStore{}
	mockIAMStore := &MockIAMStore{}
	service := NewMonitoringService(mockStore, mockIAMStore)

	orgID := uuid.New()
	start := time.Now().AddDate(0, -1, 0)
	end := time.Now()

	// Mock current usage
	currentUsage := []models.ResourceUsage{
		{
			ID:           uuid.New(),
			ResourceType: models.ResourceTypeCompute,
			ResourceID:   "compute-1",
			ResourceName: "web-server",
			Status:       models.ResourceStatusActive,
			CostPerHour:  0.05,
			TotalCost:    12.50,
			Currency:     "USD",
		},
	}

	// Mock historical usage
	historicalUsage := []models.ResourceUsageHistory{
		{
			ResourceType: models.ResourceTypeCompute,
			TotalCost:    15.00,
		},
	}

	// Mock summary
	expectedSummary := &models.ResourceUsageSummary{
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

	// Mock billing period
	billingPeriod := &models.BillingPeriod{
		OrganizationID: orgID,
		PeriodStart:    start,
		PeriodEnd:      end,
		TotalCost:      25.50,
		Currency:       "USD",
		IsPaid:         false,
	}

	mockStore.On("GetResourceUsageByOrganization", orgID.String()).Return(currentUsage, nil)
	mockStore.On("GetResourceUsageHistoryByOrganization", orgID.String(), start, end).Return(historicalUsage, nil)
	mockStore.On("GetResourceUsageSummary", orgID.String(), start, end).Return(expectedSummary, nil)
	mockStore.On("GetCurrentBillingPeriod", orgID.String()).Return(billingPeriod, nil)

	result, err := service.GetResourceUsageSummary(orgID.String(), start, end)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedSummary, &result.Summary)
	assert.Equal(t, billingPeriod, result.BillingPeriod)
	assert.Equal(t, historicalUsage, result.HistoricalUsage)
	mockStore.AssertExpectations(t)
}

func TestGenerateMonthlyBilling(t *testing.T) {
	mockStore := &MockMonitoringStore{}
	mockIAMStore := &MockIAMStore{}
	service := NewMonitoringService(mockStore, mockIAMStore)

	orgID := uuid.New()
	year := 2024
	month := 1

	// Mock monthly usage history data
	history := []models.ResourceUsageHistory{
		{
			ResourceType: models.ResourceTypeCompute,
			TotalCost:    15.00,
		},
		{
			ResourceType: models.ResourceTypeStorage,
			TotalCost:    5.25,
		},
		{
			ResourceType: models.ResourceTypeLambda,
			TotalCost:    3.25,
		},
		{
			ResourceType: models.ResourceTypeRDB,
			TotalCost:    2.00,
		},
	}

	mockStore.On("GetMonthlyUsageHistory", orgID.String(), year, month).Return(history, nil)
	mockStore.On("CreateBillingPeriod", mock.AnythingOfType("*models.BillingPeriod")).Return(nil)

	result, err := service.GenerateMonthlyBilling(orgID.String(), year, month)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 25.50, result.TotalCost) // 15.00 + 5.25 + 3.25 + 2.00
	assert.Equal(t, "USD", result.Currency)
	assert.False(t, result.IsPaid)

	mockStore.AssertExpectations(t)
}

func TestGetActiveResources(t *testing.T) {
	mockStore := &MockMonitoringStore{}
	mockIAMStore := &MockIAMStore{}
	service := NewMonitoringService(mockStore, mockIAMStore)

	orgID := uuid.New()

	expectedResources := []models.ResourceUsage{
		{
			ID:           uuid.New(),
			ResourceType: models.ResourceTypeCompute,
			ResourceID:   "compute-1",
			ResourceName: "web-server",
			Status:       models.ResourceStatusActive,
			CPUUsage:     45.2,
			MemoryUsage:  1024,
			StorageUsage: 512,
			NetworkUsage: 256,
			CostPerHour:  0.05,
			TotalCost:    12.50,
			Currency:     "USD",
		},
		{
			ID:           uuid.New(),
			ResourceType: models.ResourceTypeStorage,
			ResourceID:   "storage-1",
			ResourceName: "data-bucket",
			Status:       models.ResourceStatusActive,
			CPUUsage:     0,
			MemoryUsage:  0,
			StorageUsage: 2048,
			NetworkUsage: 128,
			CostPerHour:  0.01,
			TotalCost:    8.75,
			Currency:     "USD",
		},
	}

	mockStore.On("GetActiveResourcesByOrganization", orgID.String()).Return(expectedResources, nil)

	result, err := service.GetActiveResources(orgID.String())

	assert.NoError(t, err)
	assert.Equal(t, expectedResources, result)
	mockStore.AssertExpectations(t)
}

func TestGetMonthlyUsageTrends(t *testing.T) {
	mockStore := &MockMonitoringStore{}
	mockIAMStore := &MockIAMStore{}
	service := NewMonitoringService(mockStore, mockIAMStore)

	orgID := uuid.New()
	months := 6

	expectedTrends := []models.MonthlyUsageTrend{
		{Month: "2024-01", TotalCost: 45.25, Resources: 8},
		{Month: "2023-12", TotalCost: 38.50, Resources: 7},
		{Month: "2023-11", TotalCost: 42.75, Resources: 6},
	}

	mockStore.On("GetMonthlyUsageTrends", orgID.String(), months).Return(expectedTrends, nil)

	result, err := service.GetMonthlyUsageTrends(orgID.String(), months)

	assert.NoError(t, err)
	assert.Equal(t, expectedTrends, result)
	mockStore.AssertExpectations(t)
}

func TestUpdateMonitoringMetrics(t *testing.T) {
	mockStore := &MockMonitoringStore{}
	mockIAMStore := &MockIAMStore{}
	service := NewMonitoringService(mockStore, mockIAMStore)

	orgID := uuid.New()
	resourceID := "test-resource-123"
	resourceType := models.ResourceTypeCompute
	metrics := map[string]interface{}{
		"cpu_usage":     45.2,
		"memory_usage":  1024,
		"storage_usage": 512,
		"network_usage": 256,
	}

	// Mock existing monitoring metrics
	existingMetrics := &models.MonitoringMetrics{
		ID:             uuid.New(),
		ResourceID:     resourceID,
		ResourceType:   resourceType,
		CPUUsage:       30.0,
		MemoryUsage:    512,
		StorageUsage:   256,
		NetworkUsage:   128,
		OrganizationID: orgID,
	}

	mockStore.On("GetMonitoringMetricsByResource", resourceID, resourceType).Return(existingMetrics, nil)
	mockStore.On("UpdateMonitoringMetrics", mock.AnythingOfType("*models.MonitoringMetrics")).Return(nil)

	err := service.UpdateMonitoringMetrics(orgID, resourceID, resourceType, metrics)

	assert.NoError(t, err)
	mockStore.AssertExpectations(t)
}
