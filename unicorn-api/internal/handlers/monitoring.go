package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"unicorn-api/internal/auth"
	"unicorn-api/internal/common/errors"
	"unicorn-api/internal/common/validation"
	"unicorn-api/internal/config"
	"unicorn-api/internal/middleware"
	"unicorn-api/internal/models"
	"unicorn-api/internal/services"
	"unicorn-api/internal/stores"
)

// MonitoringHandler handles monitoring-related operations
type MonitoringHandler struct {
	service         *services.MonitoringService
	validator       *validation.Validator
	config          *config.Config
	iamStore        stores.IAMStore
	monitoringStore stores.MonitoringStore
}

// NewMonitoringHandler creates a new monitoring handler
func NewMonitoringHandler(cfg *config.Config, iamStore stores.IAMStore, monitoringStore stores.MonitoringStore) *MonitoringHandler {
	return &MonitoringHandler{
		service:         services.NewMonitoringService(monitoringStore, iamStore),
		validator:       validation.NewValidator(),
		config:          cfg,
		iamStore:        iamStore,
		monitoringStore: monitoringStore,
	}
}

// getClaims extracts claims from the request
func (h *MonitoringHandler) getClaims(c *gin.Context) (*auth.Claims, error) {
	claims, exists := middleware.GetClaimsFromContext(c)
	if !exists {
		return nil, errors.ErrUnauthorized
	}
	return claims, nil
}

// hasPermission checks if user has the required permission
func (h *MonitoringHandler) hasPermission(claims *auth.Claims, resource string, perm int) bool {
	role, err := h.iamStore.GetRoleByID(claims.RoleID)
	if err != nil {
		return false
	}
	for _, p := range role.Permissions {
		if int(p) == 0 { // Read permission
			return true
		}
	}
	return false
}

// GetResourceUsage godoc
// @Summary Get resource usage summary
// @Description Get a summary of resource usage for the current organization
// @Tags Monitoring
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start query string false "Start date (YYYY-MM-DD)"
// @Param end query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} models.ResourceUsageResponse
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/monitoring/usage [get]
func (h *MonitoringHandler) GetResourceUsage(c *gin.Context) {
	claims, err := h.getClaims(c)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	if !h.hasPermission(claims, "monitoring", 0) {
		errors.RespondWithPermissionError(c, "view resource usage")
		return
	}

	// Parse date parameters
	startStr := c.Query("start")
	endStr := c.Query("end")

	var start, end time.Time
	var err2 error

	if startStr == "" {
		// Default to 30 days ago
		start = time.Now().AddDate(0, 0, -30)
	} else {
		start, err2 = time.Parse("2006-01-02", startStr)
		if err2 != nil {
			errors.RespondWithError(c, errors.ErrBadRequest.WithDetails("Invalid start date format"))
			return
		}
	}

	if endStr == "" {
		// Default to now
		end = time.Now()
	} else {
		end, err2 = time.Parse("2006-01-02", endStr)
		if err2 != nil {
			errors.RespondWithError(c, errors.ErrBadRequest.WithDetails("Invalid end date format"))
			return
		}
	}

	// Get account to get organization ID
	account, err := h.iamStore.GetAccountByID(claims.AccountID)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	usage, err := h.service.GetResourceUsageSummary(account.OrganizationID.String(), start, end)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, usage)
}

// GetMonitoringMetrics godoc
// @Summary Get monitoring metrics for a resource
// @Description Get real-time monitoring metrics for a specific resource
// @Tags Monitoring
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resource_id path string true "Resource ID"
// @Param resource_type path string true "Resource Type"
// @Success 200 {object} models.MonitoringMetrics
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/monitoring/metrics/{resource_type}/{resource_id} [get]
func (h *MonitoringHandler) GetMonitoringMetrics(c *gin.Context) {
	claims, err := h.getClaims(c)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	if !h.hasPermission(claims, "monitoring", 0) {
		errors.RespondWithPermissionError(c, "view monitoring metrics")
		return
	}

	resourceID := c.Param("resource_id")
	resourceTypeStr := c.Param("resource_type")

	resourceType := models.ResourceType(resourceTypeStr)
	if !isValidResourceType(resourceType) {
		errors.RespondWithError(c, errors.ErrBadRequest.WithDetails("Invalid resource type"))
		return
	}

	metrics, err := h.service.GetMonitoringMetrics(resourceID, resourceType)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// UpdateMonitoringMetrics godoc
// @Summary Update monitoring metrics for a resource
// @Description Update real-time monitoring metrics for a specific resource
// @Tags Monitoring
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resource_id path string true "Resource ID"
// @Param resource_type path string true "Resource Type"
// @Param metrics body map[string]interface{} true "Monitoring metrics"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/monitoring/metrics/{resource_type}/{resource_id} [put]
func (h *MonitoringHandler) UpdateMonitoringMetrics(c *gin.Context) {
	claims, err := h.getClaims(c)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	if !h.hasPermission(claims, "monitoring", 1) {
		errors.RespondWithPermissionError(c, "update monitoring metrics")
		return
	}

	resourceID := c.Param("resource_id")
	resourceTypeStr := c.Param("resource_type")

	resourceType := models.ResourceType(resourceTypeStr)
	if !isValidResourceType(resourceType) {
		errors.RespondWithError(c, errors.ErrBadRequest.WithDetails("Invalid resource type"))
		return
	}

	var metrics map[string]interface{}
	if err := c.ShouldBindJSON(&metrics); err != nil {
		errors.RespondWithError(c, errors.ErrBadRequest.WithDetails("Invalid metrics format"))
		return
	}

	// Get account to get organization ID
	account, err := h.iamStore.GetAccountByID(claims.AccountID)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	orgID, err := uuid.Parse(account.OrganizationID.String())
	if err != nil {
		errors.RespondWithError(c, errors.ErrUnauthorized.WithDetails("Invalid organization ID"))
		return
	}

	err = h.service.UpdateMonitoringMetrics(orgID, resourceID, resourceType, metrics)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Metrics updated successfully"})
}

// GetBillingHistory godoc
// @Summary Get billing history
// @Description Get billing history for the current organization
// @Tags Monitoring
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.BillingPeriod
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/monitoring/billing [get]
func (h *MonitoringHandler) GetBillingHistory(c *gin.Context) {
	claims, err := h.getClaims(c)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	if !h.hasPermission(claims, "monitoring", 0) {
		errors.RespondWithPermissionError(c, "view billing history")
		return
	}

	// Get account to get organization ID
	account, err := h.iamStore.GetAccountByID(claims.AccountID)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	billingHistory, err := h.service.GetBillingHistory(account.OrganizationID.String())
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, billingHistory)
}

// GenerateMonthlyBilling godoc
// @Summary Generate monthly billing
// @Description Generate monthly billing for a specific month
// @Tags Monitoring
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param year query int true "Year"
// @Param month query int true "Month (1-12)"
// @Success 200 {object} models.BillingPeriod
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/monitoring/billing/generate [post]
func (h *MonitoringHandler) GenerateMonthlyBilling(c *gin.Context) {
	claims, err := h.getClaims(c)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	if !h.hasPermission(claims, "monitoring", 1) {
		errors.RespondWithPermissionError(c, "generate billing")
		return
	}

	yearStr := c.Query("year")
	monthStr := c.Query("month")

	if yearStr == "" || monthStr == "" {
		errors.RespondWithError(c, errors.ErrBadRequest.WithDetails("Year and month are required"))
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		errors.RespondWithError(c, errors.ErrBadRequest.WithDetails("Invalid year format"))
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		errors.RespondWithError(c, errors.ErrBadRequest.WithDetails("Invalid month format"))
		return
	}

	// Get account to get organization ID
	account, err := h.iamStore.GetAccountByID(claims.AccountID)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	billingPeriod, err := h.service.GenerateMonthlyBilling(account.OrganizationID.String(), year, month)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, billingPeriod)
}

// GetMonthlyUsageTrends godoc
// @Summary Get monthly usage trends
// @Description Get monthly usage trends for the current organization
// @Tags Monitoring
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param months query int false "Number of months (default: 6)"
// @Success 200 {array} models.MonthlyUsageTrend
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/monitoring/trends [get]
func (h *MonitoringHandler) GetMonthlyUsageTrends(c *gin.Context) {
	claims, err := h.getClaims(c)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	if !h.hasPermission(claims, "monitoring", 0) {
		errors.RespondWithPermissionError(c, "view usage trends")
		return
	}

	monthsStr := c.DefaultQuery("months", "6")
	months, err := strconv.Atoi(monthsStr)
	if err != nil || months < 1 || months > 24 {
		errors.RespondWithError(c, errors.ErrBadRequest.WithDetails("Invalid months parameter"))
		return
	}

	// Get account to get organization ID
	account, err := h.iamStore.GetAccountByID(claims.AccountID)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	trends, err := h.service.GetMonthlyUsageTrends(account.OrganizationID.String(), months)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, trends)
}

// GetActiveResources godoc
// @Summary Get active resources
// @Description Get all active resources for the current organization
// @Tags Monitoring
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.ResourceUsage
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/monitoring/resources/active [get]
func (h *MonitoringHandler) GetActiveResources(c *gin.Context) {
	claims, err := h.getClaims(c)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	if !h.hasPermission(claims, "monitoring", 0) {
		errors.RespondWithPermissionError(c, "view active resources")
		return
	}

	// Get account to get organization ID
	account, err := h.iamStore.GetAccountByID(claims.AccountID)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	activeResources, err := h.service.GetActiveResources(account.OrganizationID.String())
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, activeResources)
}

// TrackResourceCreation godoc
// @Summary Track resource creation
// @Description Track when a new resource is created
// @Tags Monitoring
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ResourceUsageRequest true "Resource usage request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/monitoring/track/create [post]
func (h *MonitoringHandler) TrackResourceCreation(c *gin.Context) {
	claims, err := h.getClaims(c)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	if !h.hasPermission(claims, "monitoring", 1) {
		errors.RespondWithPermissionError(c, "track resource creation")
		return
	}

	var req models.ResourceUsageRequest
	if err := h.validator.BindAndValidate(c, &req); err != nil {
		errors.RespondWithError(c, err)
		return
	}

	accountID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		errors.RespondWithError(c, errors.ErrUnauthorized.WithDetails("Invalid account ID"))
		return
	}

	// Get account to get organization ID
	account, err := h.iamStore.GetAccountByID(claims.AccountID)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	orgID, err := uuid.Parse(account.OrganizationID.String())
	if err != nil {
		errors.RespondWithError(c, errors.ErrUnauthorized.WithDetails("Invalid organization ID"))
		return
	}

	configuration := ""
	if req.Configuration != nil {
		configuration = *req.Configuration
	}

	err = h.service.TrackResourceCreation(accountID, orgID, req.ResourceType, req.ResourceID, req.ResourceName, configuration)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resource creation tracked successfully"})
}

// TrackResourceUpdate godoc
// @Summary Track resource update
// @Description Track when a resource is updated
// @Tags Monitoring
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resource_id path string true "Resource ID"
// @Param resource_type path string true "Resource Type"
// @Param request body map[string]interface{} true "Update request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/monitoring/track/{resource_type}/{resource_id} [put]
func (h *MonitoringHandler) TrackResourceUpdate(c *gin.Context) {
	claims, err := h.getClaims(c)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	if !h.hasPermission(claims, "monitoring", 1) {
		errors.RespondWithPermissionError(c, "track resource updates")
		return
	}

	resourceID := c.Param("resource_id")
	resourceTypeStr := c.Param("resource_type")

	resourceType := models.ResourceType(resourceTypeStr)
	if !isValidResourceType(resourceType) {
		errors.RespondWithError(c, errors.ErrBadRequest.WithDetails("Invalid resource type"))
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		errors.RespondWithError(c, errors.ErrBadRequest.WithDetails("Invalid update data format"))
		return
	}

	status := models.ResourceStatusActive
	if statusStr, ok := updateData["status"].(string); ok {
		status = models.ResourceStatus(statusStr)
	}

	err = h.service.TrackResourceUpdate(resourceID, resourceType, status, updateData)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resource update tracked successfully"})
}

// TrackResourceDeletion godoc
// @Summary Track resource deletion
// @Description Track when a resource is deleted
// @Tags Monitoring
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resource_id path string true "Resource ID"
// @Param resource_type path string true "Resource Type"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/monitoring/track/{resource_type}/{resource_id} [delete]
func (h *MonitoringHandler) TrackResourceDeletion(c *gin.Context) {
	claims, err := h.getClaims(c)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	if !h.hasPermission(claims, "monitoring", 1) {
		errors.RespondWithPermissionError(c, "track resource deletion")
		return
	}

	resourceID := c.Param("resource_id")
	resourceTypeStr := c.Param("resource_type")

	resourceType := models.ResourceType(resourceTypeStr)
	if !isValidResourceType(resourceType) {
		errors.RespondWithError(c, errors.ErrBadRequest.WithDetails("Invalid resource type"))
		return
	}

	err = h.service.TrackResourceDeletion(resourceID, resourceType)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resource deletion tracked successfully"})
}

// isValidResourceType checks if the resource type is valid
func isValidResourceType(resourceType models.ResourceType) bool {
	validTypes := []models.ResourceType{
		models.ResourceTypeCompute,
		models.ResourceTypeLambda,
		models.ResourceTypeStorage,
		models.ResourceTypeRDB,
		models.ResourceTypeSecret,
	}

	for _, validType := range validTypes {
		if resourceType == validType {
			return true
		}
	}
	return false
}
