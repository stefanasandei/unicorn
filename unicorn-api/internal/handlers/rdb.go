package handlers

import (
	"net/http"

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

// RDBHandler handles RDB-related operations
type RDBHandler struct {
	service   *services.RDBService
	validator *validation.Validator
	config    *config.Config
	iamStore  stores.IAMStore
}

// NewRDBHandler creates a new RDB handler
func NewRDBHandler(cfg *config.Config, iamStore stores.IAMStore) *RDBHandler {
	return &RDBHandler{
		service:   services.NewRDBService(),
		validator: validation.NewValidator(),
		config:    cfg,
		iamStore:  iamStore,
	}
}

// getClaims extracts claims from the request
func (h *RDBHandler) getClaims(c *gin.Context) (*auth.Claims, error) {
	claims, exists := middleware.GetClaimsFromContext(c)
	if !exists {
		return nil, errors.ErrUnauthorized
	}
	return claims, nil
}

// hasPermission checks if user has the required permission
func (h *RDBHandler) hasPermission(claims *auth.Claims, resource string, perm int) bool {
	role, err := h.iamStore.GetRoleByID(claims.RoleID)
	if err != nil {
		return false
	}
	for _, p := range role.Permissions {
		if int(p) == perm {
			return true
		}
	}
	return false
}

// CreateRDB godoc
// @Summary Create an RDB instance
// @Description Create a new database instance with the specified type and configuration
// @Tags RDB
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.RDBCreateRequest true "RDB instance creation request"
// @Success 200 {object} models.RDBInstanceInfo
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/rdb/create [post]
func (h *RDBHandler) CreateRDB(c *gin.Context) {
	claims, err := h.getClaims(c)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	if !h.hasPermission(claims, "rdb", 1) {
		errors.RespondWithPermissionError(c, "create RDB instances")
		return
	}

	var req models.RDBCreateRequest
	if err := h.validator.BindAndValidate(c, &req); err != nil {
		errors.RespondWithError(c, err)
		return
	}

	userID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		errors.RespondWithError(c, errors.ErrUnauthorized.WithDetails("Invalid user ID in token"))
		return
	}

	instanceInfo, err := h.service.CreateRDBInstance(userID, req)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, instanceInfo)
}

// ListRDB godoc
// @Summary List RDB instances
// @Description List all database instances owned by the authenticated user
// @Tags RDB
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.RDBInstanceInfo
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/rdb/list [get]
func (h *RDBHandler) ListRDB(c *gin.Context) {
	claims, err := h.getClaims(c)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	if !h.hasPermission(claims, "rdb", 0) {
		errors.RespondWithPermissionError(c, "list RDB instances")
		return
	}

	userID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		errors.RespondWithError(c, errors.ErrUnauthorized.WithDetails("Invalid user ID in token"))
		return
	}

	instances, err := h.service.ListRDBInstances(userID)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, instances)
}

// DeleteRDB godoc
// @Summary Delete an RDB instance
// @Description Delete a database instance by ID
// @Tags RDB
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Instance ID"
// @Success 204 {string} string "Instance deleted successfully"
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/rdb/{id} [delete]
func (h *RDBHandler) DeleteRDB(c *gin.Context) {
	claims, err := h.getClaims(c)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	if !h.hasPermission(claims, "rdb", 2) {
		errors.RespondWithPermissionError(c, "delete RDB instances")
		return
	}

	instanceID := c.Param("id")
	if instanceID == "" {
		errors.RespondWithError(c, errors.ErrBadRequest.WithDetails("Instance ID is required"))
		return
	}

	userID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		errors.RespondWithError(c, errors.ErrUnauthorized.WithDetails("Invalid user ID in token"))
		return
	}

	if err := h.service.DeleteRDBInstance(userID, instanceID); err != nil {
		errors.RespondWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
} 