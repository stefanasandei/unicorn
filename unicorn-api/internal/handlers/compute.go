package handlers

import (
	"fmt"
	"log"
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
	service           *services.ComputeService
	monitoringService *services.MonitoringService
	validator         *validation.Validator
	config            *config.Config
	iamStore          stores.IAMStore
}

// NewComputeHandler creates a new compute handler
func NewComputeHandler(cfg *config.Config, iamStore stores.IAMStore, monitoringService *services.MonitoringService) *ComputeHandler {
	return &ComputeHandler{
		service:           services.NewComputeService(),
		monitoringService: monitoringService,
		validator:         validation.NewValidator(),
		config:            cfg,
		iamStore:          iamStore,
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

	// Track resource creation for monitoring
	if h.monitoringService != nil {
		// Get account to get organization ID
		account, err := h.iamStore.GetAccountByID(claims.AccountID)
		if err == nil {
			configuration := fmt.Sprintf(`{"image":"%s","preset":"%s"}`, req.Image, req.Preset)
			err = h.monitoringService.TrackResourceCreation(
				userID,
				account.OrganizationID,
				models.ResourceTypeCompute,
				containerInfo.ID,
				req.Name,
				configuration,
			)
			if err != nil {
				// Log the error but don't fail the request
				log.Printf("Failed to track resource creation: %v", err)
			}
		}
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

// DeleteCompute godoc
// @Summary Delete a compute container
// @Description Delete a compute container by ID
// @Tags Compute
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Container ID"
// @Success 204 {string} string "Container deleted successfully"
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/compute/{id} [delete]
func (h *ComputeHandler) DeleteCompute(c *gin.Context) {
	claims, err := h.getClaims(c)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	if !h.hasPermission(claims, "compute", 2) {
		errors.RespondWithPermissionError(c, "delete compute containers")
		return
	}

	containerID := c.Param("id")
	if containerID == "" {
		errors.RespondWithError(c, errors.ErrBadRequest.WithDetails("Container ID is required"))
		return
	}

	userID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		errors.RespondWithError(c, errors.ErrUnauthorized.WithDetails("Invalid user ID in token"))
		return
	}

	if err := h.service.DeleteContainer(userID, containerID); err != nil {
		errors.RespondWithError(c, err)
		return
	}

	// Track resource deletion for monitoring
	if h.monitoringService != nil {
		err = h.monitoringService.TrackResourceDeletion(containerID, models.ResourceTypeCompute)
		if err != nil {
			// Log the error but don't fail the request
			log.Printf("Failed to track resource deletion: %v", err)
		}
	}

	c.Status(http.StatusNoContent)
}
