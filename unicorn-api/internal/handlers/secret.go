package handlers

import (
	"net/http"
	"unicorn-api/internal/auth"
	"unicorn-api/internal/config"
	"unicorn-api/internal/middleware"
	"unicorn-api/internal/models"
	"unicorn-api/internal/stores"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SecretsHandler provides endpoints for managing user secrets.
//
// @tag.name Secrets Manager
// @tag.description Manage encrypted secrets per user.
type SecretsHandler struct {
	Store    *stores.SecretStore
	Config   *config.Config
	IAMStore stores.IAMStore
}

// Helper to extract claims from Authorization header
func (h *SecretsHandler) getClaimsFromRequest(c *gin.Context) (*auth.Claims, error) {
	claims, exists := middleware.GetClaimsFromContext(c)
	if !exists {
		return nil, models.ErrTokenInvalid
	}
	return claims, nil
}

// Helper to check if user has the required permission
func (h *SecretsHandler) hasPermission(claims *auth.Claims, requiredPerm models.Permission) bool {
	role, err := h.IAMStore.GetRoleByID(claims.RoleID)
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
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse "Forbidden - insufficient permissions"
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/secrets [get]
func (h *SecretsHandler) ListSecrets(c *gin.Context) {
	claims, err := h.getClaimsFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid or missing token"})
		return
	}

	// Check if user has read permission
	if !h.hasPermission(claims, models.Read) {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "insufficient permissions to list secrets"})
		return
	}

	userID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid user id in token"})
		return
	}
	secrets, err := h.Store.ListSecrets(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
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
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse "Forbidden - insufficient permissions"
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/secrets [post]
func (h *SecretsHandler) CreateSecret(c *gin.Context) {
	claims, err := h.getClaimsFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid or missing token"})
		return
	}

	// Check if user has write permission
	if !h.hasPermission(claims, models.Write) {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "insufficient permissions to create secrets"})
		return
	}

	userID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid user id in token"})
		return
	}
	var req models.SecretBodyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	secret, err := h.Store.CreateSecret(userID, req.Name, req.Value, req.Metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
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
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse "Forbidden - insufficient permissions"
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/secrets/{id} [get]
func (h *SecretsHandler) ReadSecret(c *gin.Context) {
	claims, err := h.getClaimsFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid or missing token"})
		return
	}

	// Check if user has read permission
	if !h.hasPermission(claims, models.Read) {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "insufficient permissions to read secrets"})
		return
	}

	userID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid user id in token"})
		return
	}
	secretID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid secret id"})
		return
	}
	secret, value, err := h.Store.GetSecret(userID, secretID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
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
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse "Forbidden - insufficient permissions"
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/secrets/{id} [put]
func (h *SecretsHandler) UpdateSecret(c *gin.Context) {
	claims, err := h.getClaimsFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid or missing token"})
		return
	}

	// Check if user has write permission
	if !h.hasPermission(claims, models.Write) {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "insufficient permissions to update secrets"})
		return
	}

	userID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid user id in token"})
		return
	}
	secretID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid secret id"})
		return
	}
	var req models.UpdateSecretBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	if err := h.Store.UpdateSecret(userID, secretID, req.Value, req.Metadata); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
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
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse "Forbidden - insufficient permissions"
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/secrets/{id} [delete]
func (h *SecretsHandler) DeleteSecret(c *gin.Context) {
	claims, err := h.getClaimsFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid or missing token"})
		return
	}

	// Check if user has delete permission
	if !h.hasPermission(claims, models.Delete) {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "insufficient permissions to delete secrets"})
		return
	}

	userID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid user id in token"})
		return
	}
	secretID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid secret id"})
		return
	}
	if err := h.Store.DeleteSecret(userID, secretID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
