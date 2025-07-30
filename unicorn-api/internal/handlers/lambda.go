package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"unicorn-api/internal/auth"
	"unicorn-api/internal/config"
	"unicorn-api/internal/middleware"
	"unicorn-api/internal/models"
	"unicorn-api/internal/stores"
)

// LambdaHandler handles Lambda function execution requests
type LambdaHandler struct {
	Config    *config.Config
	IAMStore  stores.IAMStore
	LambdaURL string
}

// NewLambdaHandler creates a new Lambda handler
func NewLambdaHandler(cfg *config.Config, iamStore stores.IAMStore) *LambdaHandler {
	lambdaURL := cfg.LambdaURL
	lambdaURL = "http://localhost:6900" // Default Lambda API URL
	return &LambdaHandler{
		Config:    cfg,
		IAMStore:  iamStore,
		LambdaURL: lambdaURL,
	}
}

// ExecuteLambda godoc
// @Summary Execute a Lambda function
// @Description Execute code in a Lambda function with specified runtime and files
// @Tags Lambda
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.LambdaExecuteRequest true "Lambda execution request"
// @Success 200 {object} models.LambdaExecuteResponse
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - insufficient permissions"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/lambda/execute [post]
func (h *LambdaHandler) ExecuteLambda(c *gin.Context) {
	var req models.LambdaExecuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := h.getClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !h.hasPermission(claims, "lambda", 1) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// Forward the request to the Lambda API
	lambdaReqBody, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal request"})
		return
	}

	lambdaReq, err := http.NewRequest("POST", h.LambdaURL+"/api/v1/execute", bytes.NewBuffer(lambdaReqBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request to Lambda API"})
		return
	}
	lambdaReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(lambdaReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to Lambda API: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	// Set the same content type as the Lambda API
	c.Header("Content-Type", resp.Header.Get("Content-Type"))

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read Lambda API response"})
		return
	}

	// Return the Lambda API response
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

// TestLambda godoc
// @Summary Test a Lambda function
// @Description Test a Lambda function with specified runtime and files
// @Tags Lambda
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.LambdaExecuteRequest true "Lambda test request"
// @Success 200 {object} models.LambdaExecuteResponse
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - insufficient permissions"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/lambda/test [post]
func (h *LambdaHandler) TestLambda(c *gin.Context) {
	var req models.LambdaExecuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := h.getClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !h.hasPermission(claims, "lambda", 0) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// Forward the request to the Lambda API
	lambdaReqBody, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal request"})
		return
	}

	lambdaReq, err := http.NewRequest("POST", h.LambdaURL+"/api/v1/test", bytes.NewBuffer(lambdaReqBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request to Lambda API"})
		return
	}
	lambdaReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(lambdaReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to Lambda API: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	// Set the same content type as the Lambda API
	c.Header("Content-Type", resp.Header.Get("Content-Type"))

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read Lambda API response"})
		return
	}

	// Return the Lambda API response
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

// Helpers
func (h *LambdaHandler) getClaims(c *gin.Context) (*auth.Claims, error) {
	claims, exists := middleware.GetClaimsFromContext(c)
	if !exists {
		return nil, fmt.Errorf("authentication required")
	}
	return claims, nil
}

func (h *LambdaHandler) hasPermission(claims *auth.Claims, resource string, perm int) bool {
	// Look up the user's role and check if it has the required permission
	role, err := h.IAMStore.GetRoleByID(claims.RoleID)
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
