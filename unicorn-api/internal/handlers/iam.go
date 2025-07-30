package handlers

import (
	"net/http"
	"strings"
	"unicorn-api/internal/auth"
	"unicorn-api/internal/config"
	"unicorn-api/internal/middleware"
	"unicorn-api/internal/models"
	"unicorn-api/internal/stores"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// --- Utility structures ---

type IAMHandler struct {
	store  stores.IAMStore
	config *config.Config
}

func NewIAMHandler(store stores.IAMStore, cfg *config.Config) *IAMHandler {
	return &IAMHandler{store: store, config: cfg}
}

// --- Request/Response structures ---

// CreateRoleRequest represents the request body for creating a role
// swagger:model
type CreateRoleRequest struct {
	// The name of the role (e.g., admin, user, moderator)
	// example: admin
	Name string `json:"name" binding:"required"`
	// The permissions assigned to the role (0=Read, 1=Write, 2=Delete)
	// example: [0,1,2]
	Permissions []models.Permission `json:"permissions" binding:"required"`
}

// CreateRoleResponse represents the response when creating a role
// swagger:model
type CreateRoleResponse struct {
	// The created role
	Role models.Role `json:"role"`
	// Success message
	// example: Role created successfully
	Message string `json:"message"`
}

// AssignRoleRequest represents the request body for assigning a role to an account
// swagger:model
type AssignRoleRequest struct {
	// The ID of the account to assign the role to
	// example: 123e4567-e89b-12d3-a456-426614174000
	AccountID string `json:"account_id" binding:"required"`
	// The ID of the role to assign
	// example: 123e4567-e89b-12d3-a456-426614174000
	RoleID string `json:"role_id" binding:"required"`
}

// AssignRoleResponse represents the response when assigning a role
// swagger:model
type AssignRoleResponse struct {
	// Success message
	// example: Role assigned successfully
	Message string `json:"message"`
	// The ID of the account
	// example: 123e4567-e89b-12d3-a456-426614174000
	AccountID string `json:"account_id"`
	// The ID of the role
	// example: 123e4567-e89b-12d3-a456-426614174000
	RoleID string `json:"role_id"`
}

// CreateOrganizationRequest represents the request body for creating an organization
// swagger:model
type CreateOrganizationRequest struct {
	// The name of the organization
	// example: Acme Corporation
	Name string `json:"name" binding:"required"`
}

// CreateOrganizationResponse represents the response when creating an organization
// swagger:model
type CreateOrganizationResponse struct {
	// The created organization
	Organization models.Organization `json:"organization"`
	// Success message
	// example: Organization created successfully
	Message string `json:"message"`
}

// CreateUserRequest represents the request body for creating a user
// swagger:model
type CreateUserRequest struct {
	// The display name of the user
	// example: John Doe
	Name string `json:"name" binding:"required"`
	// The email address of the user (must be unique)
	// example: john.doe@example.com
	Email string `json:"email" binding:"required,email"`
	// The password for the user account (will be hashed)
	// example: securePassword123
	Password string `json:"password" binding:"required,min=8"`
	// The ID of the role to assign to the user
	// example: 123e4567-e89b-12d3-a456-426614174000
	RoleID string `json:"role_id" binding:"required"`
}

// CreateUserResponse represents the response when creating a user
// swagger:model
type CreateUserResponse struct {
	// The created user account
	Account models.Account `json:"account"`
	// Success message
	// example: User created successfully
	Message string `json:"message"`
}

// LoginRequest represents the request body for user authentication
// swagger:model
type LoginRequest struct {
	// The email address of the user
	// example: john.doe@example.com
	Email string `json:"email" binding:"required,email"`
	// The password for the user account
	// example: securePassword123
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the response when a user logs in
// swagger:model
type LoginResponse struct {
	// The JWT token for authentication
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoiMTIzZTQ1NjctZTg5Yi0xMmQzLWE0NTYtNDI2NjE0MTc0MDAwIiwicm9sZV9pZCI6IjEyM2U0NTY3LWU4OWItMTJkMy1hNDU2LTQyNjYxNDE3NDAwMCIsImV4cCI6MTcwNDE2ODAwMH0.example_signature
	Token string `json:"token"`
	// The type of token
	// example: Bearer
	TokenType string `json:"token_type"`
	// The expiration time of the token
	// example: 2024-01-01T12:00:00Z
	ExpiresAt time.Time `json:"expires_at"`
	// Success message
	// example: Login successful
	Message string `json:"message"`
}

// RefreshTokenRequest represents the request body for refreshing a JWT token
// swagger:model
type RefreshTokenRequest struct {
	// The JWT token to refresh
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoiMTIzZTQ1NjctZTg5Yi0xMmQzLWE0NTYtNDI2NjE0MTc0MDAwIiwicm9zZV9pZCI6IjEyM2U0NTY3LWU4OWItMTJkMy1hNDU2LTQyNjYxNDE3NDAwMCIsImV4cCI6MTcwNDE2ODAwMH0.example_signature
	Token string `json:"token" binding:"required"`
}

// RefreshTokenResponse represents the response when refreshing a token
// swagger:model
type RefreshTokenResponse struct {
	// The new JWT token
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoiMTIzZTQ1NjctZTg5Yi0xMmQzLWE0NTYtNDI2NjE0MTc0MDAwIiwicm9zZV9pZCI6IjEyM2U0NTY3LWU4OWItMTJkMy1hNDU2LTQyNjYxNDE3NDAwMCIsImV4cCI6MTcwNDE2ODAwMH0.new_signature
	Token string `json:"token"`
	// The type of token
	// example: Bearer
	TokenType string `json:"token_type"`
	// The expiration time of the token
	// example: 2024-01-01T12:00:00Z
	ExpiresAt time.Time `json:"expires_at"`
	// Success message
	// example: Token refreshed successfully
	Message string `json:"message"`
}

// TokenClaimsResponse represents the claims returned in the validate token endpoint
// swagger:model
type TokenClaimsResponse struct {
	// The account ID
	// example: 123e4567-e89b-12d3-a456-426614174000
	AccountID string `json:"account_id"`
	// The role ID
	// example: 123e4567-e89b-12d3-a456-426614174000
	RoleID string `json:"role_id"`
	// The token expiration time
	// example: 2024-01-01T12:00:00Z
	ExpiresAt time.Time `json:"expires_at"`
}

// ValidateTokenResponse represents the response when validating a token
// swagger:model
type ValidateTokenResponse struct {
	// The token claims
	Claims TokenClaimsResponse `json:"claims"`
	// Whether the token is valid
	// example: true
	Valid bool `json:"valid"`
	// Success message
	// example: Token is valid
	Message string `json:"message"`
}

// ErrorResponse represents a standard error response
// swagger:model
type ErrorResponse struct {
	// The error message
	// example: Invalid request
	Error string `json:"error"`
	// Additional error details
	// example: Field 'email' is required
	Details string `json:"details,omitempty"`
	// The HTTP status code
	// example: 400
	StatusCode int `json:"status_code"`
	// The timestamp when the error occurred
	// example: 2024-01-01T12:00:00Z
	Timestamp time.Time `json:"timestamp"`
}

// --- High-Level Handler functions ---

// CreateRole godoc
// @Summary      Create a new role
// @Description  Create a new role with specified permissions. Permissions are: 0=Read, 1=Write, 2=Delete
// @Tags         IAM
// @Accept       json
// @Produce      json
// @Param        role  body  CreateRoleRequest  true  "Role information"
// @Success      201   {object}  CreateRoleResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/v1/roles [post]
func (h *IAMHandler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:      "Invalid request",
			Details:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Timestamp:  time.Now(),
		})
		return
	}

	role := &models.Role{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Name:        req.Name,
		Permissions: req.Permissions,
	}
	if err := h.store.CreateRole(role); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
			Timestamp:  time.Now(),
		})
		return
	}

	c.JSON(http.StatusCreated, CreateRoleResponse{
		Role:    *role,
		Message: "Role created successfully",
	})
}

// AssignRole godoc
// @Summary      Assign a role to an account
// @Description  Assign a role to a user or bot account
// @Tags         IAM
// @Accept       json
// @Produce      json
// @Param        assignment  body  AssignRoleRequest  true  "Role assignment information"
// @Success      200   {object}  AssignRoleResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/v1/roles/assign [post]
func (h *IAMHandler) AssignRole(c *gin.Context) {
	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:      "Invalid request",
			Details:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Timestamp:  time.Now(),
		})
		return
	}
	if err := h.store.AssignRole(req.AccountID, req.RoleID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
			Timestamp:  time.Now(),
		})
		return
	}
	c.JSON(http.StatusOK, AssignRoleResponse{
		Message:   "Role assigned successfully",
		AccountID: req.AccountID,
		RoleID:    req.RoleID,
	})
}

// CreateOrganization godoc
// @Summary      Create a new organization
// @Description  Create a new organization that can contain multiple accounts
// @Tags         IAM
// @Accept       json
// @Produce      json
// @Param        organization  body  CreateOrganizationRequest  true  "Organization information"
// @Success      201   {object}  CreateOrganizationResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/v1/organizations [post]
func (h *IAMHandler) CreateOrganization(c *gin.Context) {
	var req CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:      "Invalid request",
			Details:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Timestamp:  time.Now(),
		})
		return
	}
	org := &models.Organization{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      req.Name,
	}
	if err := h.store.CreateOrganization(org); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
			Timestamp:  time.Now(),
		})
		return
	}
	c.JSON(http.StatusCreated, CreateOrganizationResponse{
		Organization: *org,
		Message:      "Organization created successfully",
	})
}

// --- User-related Handler functions ---

// CreateUserInOrg godoc
// @Summary      Create a user in an organization
// @Description  Create a user account in a specific organization with the specified role
// @Tags         IAM
// @Accept       json
// @Produce      json
// @Param        org_id  path  string  true  "Organization ID"
// @Param        user  body  CreateUserRequest  true  "User information"
// @Success      201   {object}  CreateUserResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/v1/organizations/{org_id}/users [post]
func (h *IAMHandler) CreateUserInOrg(c *gin.Context) {
	orgID := c.Param("org_id")

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:      "Invalid request",
			Details:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Timestamp:  time.Now(),
		})
		return
	}

	hashed, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:      "Failed to hash password",
			StatusCode: http.StatusInternalServerError,
			Timestamp:  time.Now(),
		})
		return
	}

	account := &models.Account{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Name:           req.Name,
		Email:          req.Email,
		Type:           models.AccountTypeUser,
		PasswordHash:   hashed,
		OrganizationID: uuid.MustParse(orgID),
		RoleID:         uuid.MustParse(req.RoleID),
	}
	if err := h.store.CreateAccount(account); err != nil {
		if err == models.ErrDuplicateEmail {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:      "Email already exists",
				StatusCode: http.StatusInternalServerError,
				Timestamp:  time.Now(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
			Timestamp:  time.Now(),
		})
		return
	}
	c.JSON(http.StatusCreated, CreateUserResponse{
		Account: *account,
		Message: "User created successfully",
	})
}

// Login godoc
// @Summary      Login
// @Description  Authenticate a user and return a JWT token
// @Tags         IAM
// @Accept       json
// @Produce      json
// @Param        credentials  body  LoginRequest  true  "Login credentials"
// @Success      200   {object}  LoginResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      401   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/v1/login [post]
func (h *IAMHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:      "Invalid request",
			Details:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Timestamp:  time.Now(),
		})
		return
	}
	account, err := h.store.GetAccountByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:      "Invalid credentials",
			StatusCode: http.StatusUnauthorized,
			Timestamp:  time.Now(),
		})
		return
	}
	if !auth.CheckPasswordHash(req.Password, account.PasswordHash) {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:      "Invalid credentials",
			StatusCode: http.StatusUnauthorized,
			Timestamp:  time.Now(),
		})
		return
	}
	token, err := auth.GenerateToken(account.ID.String(), account.RoleID.String(), h.config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:      "Failed to generate token",
			StatusCode: http.StatusInternalServerError,
			Timestamp:  time.Now(),
		})
		return
	}

	// Calculate expiration time (assuming 24 hours)
	expiresAt := time.Now().Add(24 * time.Hour)

	c.JSON(http.StatusOK, LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresAt: expiresAt,
		Message:   "Login successful",
	})
}

// RefreshToken godoc
// @Summary      Refresh JWT token
// @Description  Refresh an expired JWT token
// @Tags         IAM
// @Accept       json
// @Produce      json
// @Param        token  body  RefreshTokenRequest  true  "Refresh token request"
// @Success      200   {object}  RefreshTokenResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      401   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/v1/token/refresh [post]
func (h *IAMHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:      "Invalid request",
			Details:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Timestamp:  time.Now(),
		})
		return
	}
	claims, err := auth.ValidateToken(req.Token, h.config)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:      "Invalid token",
			StatusCode: http.StatusUnauthorized,
			Timestamp:  time.Now(),
		})
		return
	}
	token, err := auth.GenerateToken(claims.AccountID, claims.RoleID, h.config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:      "Failed to generate token",
			StatusCode: http.StatusInternalServerError,
			Timestamp:  time.Now(),
		})
		return
	}

	// Calculate expiration time (assuming 24 hours)
	expiresAt := time.Now().Add(24 * time.Hour)

	c.JSON(http.StatusOK, RefreshTokenResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresAt: expiresAt,
		Message:   "Token refreshed successfully",
	})
}

// ValidateToken godoc
// @Summary      Validate JWT token
// @Description  Validate a JWT token and return its claims
// @Tags         IAM
// @Accept       json
// @Produce      json
// @Param        token  query  string  true  "JWT token"
// @Success      200   {object}  ValidateTokenResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      401   {object}  ErrorResponse
// @Router       /api/v1/token/validate [get]
func (h *IAMHandler) ValidateToken(c *gin.Context) {
	token := c.Query("token")
	token = strings.TrimPrefix(token, "Bearer: ")
	claims, err := auth.ValidateToken(token, h.config)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:      "Invalid token",
			StatusCode: http.StatusUnauthorized,
			Timestamp:  time.Now(),
		})
		return
	}
	c.JSON(http.StatusOK, ValidateTokenResponse{
		Claims: TokenClaimsResponse{
			AccountID: claims.AccountID,
			RoleID:    claims.RoleID,
			ExpiresAt: claims.ExpiresAt.Time,
		},
		Valid:   true,
		Message: "Token is valid",
	})
}

// GetDebugToken godoc
// @Summary      Get a debug API token for testing
// @Description  Returns a static API token with admin permissions for testing purposes
// @Tags         IAM
// @Accept       json
// @Produce      json
// @Success      200   {object}  map[string]string
// @Failure      500   {object}  ErrorResponse
// @Router       /api/v1/debug/token [get]
func (h *IAMHandler) GetDebugToken(c *gin.Context) {
	if h.config.DebugAPIToken == "" {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:      "Debug token not available",
			StatusCode: http.StatusInternalServerError,
			Timestamp:  time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      h.config.DebugAPIToken,
		"token_type": "Bearer",
		"message":    "This token has admin permissions and should only be used for testing",
	})
}

// --- New Response Structures ---

// GetRolesResponse represents the response for listing roles
// swagger:model
type GetRolesResponse struct {
	Roles []models.Role `json:"roles"`
}

// GetOrganizationsResponse represents the response for listing organizations and their users
// swagger:model
type GetOrganizationsResponse struct {
	OrganizationName string             `json:"organization_name"`
	Users            []OrganizationUser `json:"users"`
}

type OrganizationUser struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	RoleID string `json:"role_id"`
}

// --- New Handler Methods ---

// GetRoles godoc
// @Summary      Get all roles in the user's organization
// @Description  Returns all roles (name, permissions) for the authenticated user's organization
// @Tags         IAM
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200   {object}  GetRolesResponse
// @Failure      401   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/v1/roles [get]
func (h *IAMHandler) GetRoles(c *gin.Context) {
	claims, exists := middleware.GetClaimsFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:      "Authentication required",
			StatusCode: http.StatusUnauthorized,
			Timestamp:  time.Now(),
		})
		return
	}
	
	account, err := h.store.GetAccountByID(claims.AccountID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:      "Account not found",
			StatusCode: http.StatusUnauthorized,
			Timestamp:  time.Now(),
		})
		return
	}
	roles, err := h.store.GetRolesByOrganizationID(account.OrganizationID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
			Timestamp:  time.Now(),
		})
		return
	}
	c.JSON(http.StatusOK, GetRolesResponse{Roles: roles})
}

// GetOrganizations godoc
// @Summary      Get the user's organization and its users
// @Description  Returns the organization name and all users (name, role ID) in it
// @Tags         IAM
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200   {object}  GetOrganizationsResponse
// @Failure      401   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/v1/organizations [get]
func (h *IAMHandler) GetOrganizations(c *gin.Context) {
	claims, exists := middleware.GetClaimsFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:      "Authentication required",
			StatusCode: http.StatusUnauthorized,
			Timestamp:  time.Now(),
		})
		return
	}
	
	account, err := h.store.GetAccountByID(claims.AccountID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:      "Account not found",
			StatusCode: http.StatusUnauthorized,
			Timestamp:  time.Now(),
		})
		return
	}
	org, err := h.store.GetOrganizationByID(account.OrganizationID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
			Timestamp:  time.Now(),
		})
		return
	}
	users, err := h.store.GetAccountsByOrganizationID(account.OrganizationID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
			Timestamp:  time.Now(),
		})
		return
	}
	var orgUsers []OrganizationUser
	for _, u := range users {
		orgUsers = append(orgUsers, OrganizationUser{
			ID:     u.ID.String(),
			Name:   u.Name,
			RoleID: u.RoleID.String(),
		})
	}
	c.JSON(http.StatusOK, GetOrganizationsResponse{
		OrganizationName: org.Name,
		Users:            orgUsers,
	})
}
