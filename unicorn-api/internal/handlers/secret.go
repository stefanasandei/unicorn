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

// SecretsHandler provides endpoints for managing user secrets.
//
// @tag.name Secrets Manager
// @tag.description Manage encrypted secrets per user.
type SecretsHandler struct {
	service   *services.SecretService
	validator *validation.Validator
	config    *config.Config
	iamStore  stores.IAMStore
}

// NewSecretsHandler creates a new secrets handler
func NewSecretsHandler(secretStore stores.SecretStoreInterface, iamStore stores.IAMStore, cfg *config.Config) *SecretsHandler {
	return &SecretsHandler{
		service:   services.NewSecretService(secretStore),
		validator: validation.NewValidator(),
		config:    cfg,
		iamStore:  iamStore,
	}
}

// getClaimsFromRequest extracts claims from the Authorization header
func (h *SecretsHandler) getClaimsFromRequest(c *gin.Context) (*auth.Claims, error) {
	claims, exists := middleware.GetClaimsFromContext(c)
	if !exists {
		return nil, errors.ErrUnauthorized
	}
	return claims, nil
}

// hasPermission checks if user has the required permission
func (h *SecretsHandler) hasPermission(claims *auth.Claims, requiredPerm models.Permission) bool {
	role, err := h.iamStore.GetRoleByID(claims.RoleID)
	if err != nil {
		return false
	}

	for _, perm := range role.Permissions {
		if perm == requiredPerm {
			return true
		}
	}
	return false
}

// ListSecrets returns all secrets for the authenticated user (no values)
// @Summary List secrets
// @Description List all secrets for the authenticated user (does not return secret values)
// @Tags Secrets Manager
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.SecretResponse
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/secrets [get]
func (h *SecretsHandler) ListSecrets(c *gin.Context) {
	claims, err := h.getClaimsFromRequest(c)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	// Check if user has read permission
	if !h.hasPermission(claims, models.Read) {
		errors.RespondWithPermissionError(c, "list secrets")
		return
	}

	userID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		errors.RespondWithError(c, errors.ErrUnauthorized.WithDetails("Invalid user ID in token"))
		return
	}

	secrets, err := h.service.ListSecrets(userID)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, secrets)
}

// CreateSecret creates a new secret for the user
// @Summary Create secret
// @Description Create a new secret for the authenticated user
// @Tags Secrets Manager
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param secret body models.SecretBodyRequest true "Secret details to create"
// @Success 201 {object} models.SecretResponse
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/secrets [post]
func (h *SecretsHandler) CreateSecret(c *gin.Context) {
	claims, err := h.getClaimsFromRequest(c)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	// Check if user has write permission
	if !h.hasPermission(claims, models.Write) {
		errors.RespondWithPermissionError(c, "create secrets")
		return
	}

	userID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		errors.RespondWithError(c, errors.ErrUnauthorized.WithDetails("Invalid user ID in token"))
		return
	}

	var req models.SecretBodyRequest
	if err := h.validator.BindAndValidate(c, &req); err != nil {
		errors.RespondWithError(c, err)
		return
	}

	secret, err := h.service.CreateSecret(userID, req.Name, req.Value, req.Metadata)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	resp := models.SecretResponse{
		ID:        secret.ID,
		Name:      secret.Name,
		CreatedAt: secret.CreatedAt,
		UpdatedAt: secret.UpdatedAt,
		UserID:    secret.UserID,
		Metadata:  secret.MetadataRaw,
	}

	c.JSON(http.StatusCreated, resp)
}

// ReadSecret returns the decrypted value for the user's secret
// @Summary Get secret
// @Description Get a secret and its decrypted value for the authenticated user
// @Tags Secrets Manager
// @Produce json
// @Security BearerAuth
// @Param id path string true "Secret ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Router /api/v1/secrets/{id} [get]
func (h *SecretsHandler) ReadSecret(c *gin.Context) {
	claims, err := h.getClaimsFromRequest(c)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	// Check if user has read permission
	if !h.hasPermission(claims, models.Read) {
		errors.RespondWithPermissionError(c, "read secrets")
		return
	}

	userID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		errors.RespondWithError(c, errors.ErrUnauthorized.WithDetails("Invalid user ID in token"))
		return
	}

	secretID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		errors.RespondWithError(c, errors.ErrBadRequest.WithDetails("Invalid secret ID"))
		return
	}

	secret, value, err := h.service.GetSecret(userID, secretID)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         secret.ID,
		"name":       secret.Name,
		"value":      value,
		"created_at": secret.CreatedAt,
		"updated_at": secret.UpdatedAt,
		"user_id":    secret.UserID,
		"metadata":   secret.Metadata,
	})
}

// UpdateSecret updates the value and/or metadata for the user's secret
// @Summary Update secret
// @Description Update the value and/or metadata for a secret
// @Tags Secrets Manager
// @Accept json
// @Security BearerAuth
// @Param id path string true "Secret ID"
// @Param secret body models.UpdateSecretBody true "Secret update"
// @Success 204
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/secrets/{id} [put]
func (h *SecretsHandler) UpdateSecret(c *gin.Context) {
	claims, err := h.getClaimsFromRequest(c)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	// Check if user has write permission
	if !h.hasPermission(claims, models.Write) {
		errors.RespondWithPermissionError(c, "update secrets")
		return
	}

	userID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		errors.RespondWithError(c, errors.ErrUnauthorized.WithDetails("Invalid user ID in token"))
		return
	}

	secretID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		errors.RespondWithError(c, errors.ErrBadRequest.WithDetails("Invalid secret ID"))
		return
	}

	var req models.UpdateSecretBody
	if err := h.validator.BindAndValidate(c, &req); err != nil {
		errors.RespondWithError(c, err)
		return
	}

	if err := h.service.UpdateSecret(userID, secretID, req.Value, req.Metadata); err != nil {
		errors.RespondWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteSecret deletes the user's secret
// @Summary Delete secret
// @Description Delete a secret for the authenticated user
// @Tags Secrets Manager
// @Security BearerAuth
// @Param id path string true "Secret ID"
// @Success 204
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/secrets/{id} [delete]
func (h *SecretsHandler) DeleteSecret(c *gin.Context) {
	claims, err := h.getClaimsFromRequest(c)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	// Check if user has delete permission
	if !h.hasPermission(claims, models.Delete) {
		errors.RespondWithPermissionError(c, "delete secrets")
		return
	}

	userID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		errors.RespondWithError(c, errors.ErrUnauthorized.WithDetails("Invalid user ID in token"))
		return
	}

	secretID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		errors.RespondWithError(c, errors.ErrBadRequest.WithDetails("Invalid secret ID"))
		return
	}

	if err := h.service.DeleteSecret(userID, secretID); err != nil {
		errors.RespondWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// RotateKeys rotates all keys for the authenticated user
// @Summary Rotate keys
// @Description Rotate all encryption keys for the authenticated user
// @Tags Secrets Manager
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/secrets/rotate-keys [post]
func (h *SecretsHandler) RotateKeys(c *gin.Context) {
	claims, err := h.getClaimsFromRequest(c)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	// Check if user has write permission
	if !h.hasPermission(claims, models.Write) {
		errors.RespondWithPermissionError(c, "rotate keys")
		return
	}

	userID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		errors.RespondWithError(c, errors.ErrUnauthorized.WithDetails("Invalid user ID in token"))
		return
	}

	if err := h.service.RotateKeys(userID); err != nil {
		errors.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Keys rotated successfully",
		"user_id": userID,
	})
}

// GetKeyVersions gets all key versions for the authenticated user
// @Summary Get key versions
// @Description Get all key versions for the authenticated user
// @Tags Secrets Manager
// @Security BearerAuth
// @Success 200 {array} stores.KeyVersion
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/secrets/key-versions [get]
func (h *SecretsHandler) GetKeyVersions(c *gin.Context) {
	claims, err := h.getClaimsFromRequest(c)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	// Check if user has read permission
	if !h.hasPermission(claims, models.Read) {
		errors.RespondWithPermissionError(c, "view key versions")
		return
	}

	userID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		errors.RespondWithError(c, errors.ErrUnauthorized.WithDetails("Invalid user ID in token"))
		return
	}

	versions, err := h.service.GetKeyVersions(userID)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, versions)
}
