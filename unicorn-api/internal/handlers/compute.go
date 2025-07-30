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

// ComputeHandler handles compute-related operations
type ComputeHandler struct {
	service   *services.ComputeService
	validator *validation.Validator
	config    *config.Config
	iamStore  stores.IAMStore
}

// NewComputeHandler creates a new compute handler
func NewComputeHandler(cfg *config.Config, iamStore stores.IAMStore) *ComputeHandler {
	return &ComputeHandler{
		service:   services.NewComputeService(),
		validator: validation.NewValidator(),
		config:    cfg,
		iamStore:  iamStore,
	}
}

// getClaims extracts claims from the request
func (h *ComputeHandler) getClaims(c *gin.Context) (*auth.Claims, error) {
	claims, exists := middleware.GetClaimsFromContext(c)
	if !exists {
		return nil, errors.ErrUnauthorized
	}
	return claims, nil
}

// hasPermission checks if user has the required permission
func (h *ComputeHandler) hasPermission(claims *auth.Claims, resource string, perm int) bool {
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

// CreateCompute godoc
// @Summary Create a compute container
// @Description Create a new compute container with the specified image and configuration
// @Tags Compute
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ComputeCreateRequest true "Compute container creation request"
// @Success 200 {object} models.ComputeContainerInfo
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/compute/create [post]
func (h *ComputeHandler) CreateCompute(c *gin.Context) {
	claims, err := h.getClaims(c)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	if !h.hasPermission(claims, "compute", 1) {
		errors.RespondWithPermissionError(c, "create compute containers")
		return
	}

	var req models.ComputeCreateRequest
	if err := h.validator.BindAndValidate(c, &req); err != nil {
		errors.RespondWithError(c, err)
		return
	}

	userID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		errors.RespondWithError(c, errors.ErrUnauthorized.WithDetails("Invalid user ID in token"))
		return
	}

	containerInfo, err := h.service.CreateContainer(userID, req)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, containerInfo)
}

// ListCompute godoc
// @Summary List compute containers
// @Description List all compute containers owned by the authenticated user
// @Tags Compute
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.ComputeContainerInfo
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/compute/list [get]
func (h *ComputeHandler) ListCompute(c *gin.Context) {
	claims, err := h.getClaims(c)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	if !h.hasPermission(claims, "compute", 0) {
		errors.RespondWithPermissionError(c, "list compute containers")
		return
	}

	userID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		errors.RespondWithError(c, errors.ErrUnauthorized.WithDetails("Invalid user ID in token"))
		return
	}

	containers, err := h.service.ListContainers(userID)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, containers)
}
