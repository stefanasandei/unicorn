package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"unicorn-api/internal/config"
	"unicorn-api/internal/handlers"
	"unicorn-api/internal/models"
	"unicorn-api/internal/routes"
	"unicorn-api/internal/services"
	"unicorn-api/internal/stores"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupIntegrationTest(t *testing.T) (*gin.Engine, *stores.GORMIAMStore, func()) {
	// Create a temporary database file
	dbPath := "integration_test_" + uuid.New().String() + ".db"

	store, err := stores.NewGORMIAMStore(dbPath)
	require.NoError(t, err)

	cfg := &config.Config{
		JWTSecret:       "integration-test-secret-key",
		TokenExpiration: 24 * time.Hour,
		Environment:     "test",
	}

	secretsStore, err := stores.NewSecretStore("test.db")
	if err != nil {
		panic("failed to initialize secrets store: " + err.Error())
	}

	monitoringStore, err := stores.NewGORMMonitoringStore("test.db")
	if err != nil {
		panic("failed to initialize monitoring store: " + err.Error())
	}
	monitoringService := services.NewMonitoringService(monitoringStore, store)

	handler := handlers.NewIAMHandler(store, cfg)
	storageHandler := handlers.NewStorageHandler(&stores.GORMStorageStore{}, store, cfg)
	computeHandler := handlers.NewComputeHandler(cfg, store, monitoringService)
	lambdaHandler := handlers.NewLambdaHandler(cfg, store)
	secretsHandler := handlers.NewSecretsHandler(secretsStore, store, cfg)
	rdbHandler := handlers.NewRDBHandler(cfg, store)
	monitoringHandler := handlers.NewMonitoringHandler(cfg, store, monitoringStore)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	routes.SetupRoutes(router, handler, storageHandler, computeHandler, lambdaHandler, secretsHandler, rdbHandler, monitoringHandler, cfg)

	// Return cleanup function
	cleanup := func() {
		store.DB().Migrator().DropTable(&models.Role{}, &models.Organization{}, &models.Account{})
		os.Remove(dbPath)
	}

	return router, store, cleanup
}

// TestFullIAMWorkflow tests the complete workflow of:
// 1. Creating an organization
// 2. Creating roles with different permissions
// 3. Creating users in the organization
// 4. Assigning roles to users
// 5. Testing authentication
// 6. Testing token validation
func TestFullIAMWorkflow(t *testing.T) {
	router, _, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Step 1: Create an organization
	var orgID string
	t.Run("Create Organization", func(t *testing.T) {
		orgRequest := handlers.CreateOrganizationRequest{
			Name: "Acme Corporation",
		}

		body, _ := json.Marshal(orgRequest)
		req := httptest.NewRequest("POST", "/api/v1/organizations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response handlers.CreateOrganizationResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Acme Corporation", response.Organization.Name)
		assert.NotEqual(t, uuid.Nil, response.Organization.ID)

		// Store org ID for later use
		orgID = response.Organization.ID.String()
	})

	// Step 2: Create roles with different permission levels
	var adminRoleID, userRoleID, readOnlyRoleID string

	t.Run("Create Admin Role", func(t *testing.T) {
		roleRequest := handlers.CreateRoleRequest{
			Name:        "admin",
			Permissions: []models.Permission{models.Read, models.Write, models.Delete},
		}

		body, _ := json.Marshal(roleRequest)
		req := httptest.NewRequest("POST", "/api/v1/roles", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response handlers.CreateRoleResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "admin", response.Role.Name)
		assert.Len(t, response.Role.Permissions, 3)

		adminRoleID = response.Role.ID.String()
	})

	t.Run("Create User Role", func(t *testing.T) {
		roleRequest := handlers.CreateRoleRequest{
			Name:        "user",
			Permissions: []models.Permission{models.Read, models.Write},
		}

		body, _ := json.Marshal(roleRequest)
		req := httptest.NewRequest("POST", "/api/v1/roles", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response handlers.CreateRoleResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "user", response.Role.Name)
		assert.Len(t, response.Role.Permissions, 2)

		userRoleID = response.Role.ID.String()
	})

	t.Run("Create Read-Only Role", func(t *testing.T) {
		roleRequest := handlers.CreateRoleRequest{
			Name:        "readonly",
			Permissions: []models.Permission{models.Read},
		}

		body, _ := json.Marshal(roleRequest)
		req := httptest.NewRequest("POST", "/api/v1/roles", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response handlers.CreateRoleResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "readonly", response.Role.Name)
		assert.Len(t, response.Role.Permissions, 1)

		readOnlyRoleID = response.Role.ID.String()
	})

	// Step 3: Create users in the organization
	var adminUserID, regularUserID, readOnlyUserID string

	t.Run("Create Admin User", func(t *testing.T) {
		userRequest := handlers.CreateUserRequest{
			Name:     "Admin User",
			Email:    "admin@acme.com",
			Password: "adminpass123",
			RoleID:   adminRoleID,
		}

		body, _ := json.Marshal(userRequest)
		url := "/api/v1/organizations/" + orgID + "/users"
		req := httptest.NewRequest("POST", url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response handlers.CreateUserResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Admin User", response.Account.Name)
		assert.Equal(t, "admin@acme.com", response.Account.Email)
		assert.Equal(t, models.AccountTypeUser, response.Account.Type)

		adminUserID = response.Account.ID.String()
	})

	t.Run("Create Regular User", func(t *testing.T) {
		userRequest := handlers.CreateUserRequest{
			Name:     "Regular User",
			Email:    "user@acme.com",
			Password: "userpass123",
			RoleID:   userRoleID,
		}

		body, _ := json.Marshal(userRequest)
		url := "/api/v1/organizations/" + orgID + "/users"
		req := httptest.NewRequest("POST", url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response handlers.CreateUserResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Regular User", response.Account.Name)
		assert.Equal(t, "user@acme.com", response.Account.Email)

		regularUserID = response.Account.ID.String()
	})

	t.Run("Create Read-Only User", func(t *testing.T) {
		userRequest := handlers.CreateUserRequest{
			Name:     "Read-Only User",
			Email:    "readonly@acme.com",
			Password: "readonlypass123",
			RoleID:   readOnlyRoleID,
		}

		body, _ := json.Marshal(userRequest)
		url := "/api/v1/organizations/" + orgID + "/users"
		req := httptest.NewRequest("POST", url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response handlers.CreateUserResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Read-Only User", response.Account.Name)
		assert.Equal(t, "readonly@acme.com", response.Account.Email)

		readOnlyUserID = response.Account.ID.String()
	})

	// Step 4: Test role assignment
	t.Run("Assign Different Role to User", func(t *testing.T) {
		// First, create a new role
		newRoleRequest := handlers.CreateRoleRequest{
			Name:        "moderator",
			Permissions: []models.Permission{models.Read, models.Write},
		}

		body, _ := json.Marshal(newRoleRequest)
		req := httptest.NewRequest("POST", "/api/v1/roles", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var roleResponse handlers.CreateRoleResponse
		err := json.Unmarshal(w.Body.Bytes(), &roleResponse)
		require.NoError(t, err)

		// Store the role ID for later assignment
		userRoleID = roleResponse.Role.ID.String()
	})

	// Step 5: Test authentication for all users
	var adminToken, userToken, readOnlyToken string

	t.Run("Admin User Login", func(t *testing.T) {
		loginRequest := handlers.LoginRequest{
			Email:    "admin@acme.com",
			Password: "adminpass123",
		}

		body, _ := json.Marshal(loginRequest)
		req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response handlers.LoginResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotEmpty(t, response.Token)
		assert.Equal(t, "Bearer", response.TokenType)
		assert.Equal(t, "Login successful", response.Message)

		adminToken = response.Token
	})

	t.Run("Regular User Login", func(t *testing.T) {
		loginRequest := handlers.LoginRequest{
			Email:    "user@acme.com",
			Password: "userpass123",
		}

		body, _ := json.Marshal(loginRequest)
		req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response handlers.LoginResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotEmpty(t, response.Token)
		assert.Equal(t, "Bearer", response.TokenType)

		userToken = response.Token
	})

	t.Run("Read-Only User Login", func(t *testing.T) {
		loginRequest := handlers.LoginRequest{
			Email:    "readonly@acme.com",
			Password: "readonlypass123",
		}

		body, _ := json.Marshal(loginRequest)
		req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response handlers.LoginResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotEmpty(t, response.Token)
		assert.Equal(t, "Bearer", response.TokenType)

		readOnlyToken = response.Token
	})

	// Step 6: Test token validation
	t.Run("Validate Admin Token", func(t *testing.T) {
		url := "/api/v1/token/validate?token=" + adminToken
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response handlers.ValidateTokenResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response.Valid)
		assert.Equal(t, adminUserID, response.Claims.AccountID)
		assert.Equal(t, adminRoleID, response.Claims.RoleID)
		assert.Equal(t, "Token is valid", response.Message)
	})

	t.Run("Validate User Token", func(t *testing.T) {
		url := "/api/v1/token/validate?token=" + userToken
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response handlers.ValidateTokenResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response.Valid)
		assert.Equal(t, regularUserID, response.Claims.AccountID)
		assert.Equal(t, "Token is valid", response.Message)
	})

	t.Run("Validate Read-Only Token", func(t *testing.T) {
		url := "/api/v1/token/validate?token=" + readOnlyToken
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response handlers.ValidateTokenResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response.Valid)
		assert.Equal(t, readOnlyUserID, response.Claims.AccountID)
		assert.Equal(t, readOnlyRoleID, response.Claims.RoleID)
		assert.Equal(t, "Token is valid", response.Message)
	})

	// Step 7: Test invalid authentication scenarios
	t.Run("Invalid Login - Wrong Password", func(t *testing.T) {
		loginRequest := handlers.LoginRequest{
			Email:    "admin@acme.com",
			Password: "wrongpassword",
		}

		body, _ := json.Marshal(loginRequest)
		req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid Login - Non-existent User", func(t *testing.T) {
		loginRequest := handlers.LoginRequest{
			Email:    "nonexistent@acme.com",
			Password: "anypassword",
		}

		body, _ := json.Marshal(loginRequest)
		req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid Token Validation", func(t *testing.T) {
		url := "/api/v1/token/validate?token=invalid.token.here"
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// Step 7: Test role assignment with authentication
	t.Run("Assign Different Role to User", func(t *testing.T) {
		// Now assign the user role to the regular user using admin token
		assignRequest := handlers.AssignRoleRequest{
			AccountID: regularUserID,
			RoleID:    userRoleID,
		}

		body, _ := json.Marshal(assignRequest)
		req := httptest.NewRequest("POST", "/api/v1/roles/assign", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+adminToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var assignResponse handlers.AssignRoleResponse
		err := json.Unmarshal(w.Body.Bytes(), &assignResponse)
		require.NoError(t, err)

		assert.Equal(t, regularUserID, assignResponse.AccountID)
		assert.Equal(t, userRoleID, assignResponse.RoleID)
	})

	// Step 8: Test data persistence by verifying the data is still there
	t.Run("Verify Data Persistence", func(t *testing.T) {
		// Try to login again with the same credentials
		loginRequest := handlers.LoginRequest{
			Email:    "admin@acme.com",
			Password: "adminpass123",
		}

		body, _ := json.Marshal(loginRequest)
		req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response handlers.LoginResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotEmpty(t, response.Token)
	})
}

// TestOrganizationIsolation tests that users from different organizations are properly isolated
func TestOrganizationIsolation(t *testing.T) {
	router, _, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Create two organizations
	var org1ID, org2ID string

	t.Run("Create First Organization", func(t *testing.T) {
		orgRequest := handlers.CreateOrganizationRequest{
			Name: "Organization A",
		}

		body, _ := json.Marshal(orgRequest)
		req := httptest.NewRequest("POST", "/api/v1/organizations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response handlers.CreateOrganizationResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		org1ID = response.Organization.ID.String()
	})

	t.Run("Create Second Organization", func(t *testing.T) {
		orgRequest := handlers.CreateOrganizationRequest{
			Name: "Organization B",
		}

		body, _ := json.Marshal(orgRequest)
		req := httptest.NewRequest("POST", "/api/v1/organizations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response handlers.CreateOrganizationResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		org2ID = response.Organization.ID.String()
	})

	// Create a role that can be shared
	var sharedRoleID string
	t.Run("Create Shared Role", func(t *testing.T) {
		roleRequest := handlers.CreateRoleRequest{
			Name:        "user",
			Permissions: []models.Permission{models.Read, models.Write},
		}

		body, _ := json.Marshal(roleRequest)
		req := httptest.NewRequest("POST", "/api/v1/roles", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response handlers.CreateRoleResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		sharedRoleID = response.Role.ID.String()
	})

	// Create users in different organizations with the same email (should be allowed)
	t.Run("Create User in Organization A", func(t *testing.T) {
		userRequest := handlers.CreateUserRequest{
			Name:     "User A",
			Email:    "user@example.com",
			Password: "password123",
			RoleID:   sharedRoleID,
		}

		body, _ := json.Marshal(userRequest)
		url := "/api/v1/organizations/" + org1ID + "/users"
		req := httptest.NewRequest("POST", url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Create User in Organization B", func(t *testing.T) {
		userRequest := handlers.CreateUserRequest{
			Name:     "User B",
			Email:    "user@example.com", // Same email, different org
			Password: "password123",
			RoleID:   sharedRoleID,
		}

		body, _ := json.Marshal(userRequest)
		url := "/api/v1/organizations/" + org2ID + "/users"
		req := httptest.NewRequest("POST", url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// This should fail because email uniqueness is enforced globally
		// (This test verifies the current behavior - if you want to allow same email across orgs,
		// you'd need to modify the database schema and validation)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

// TestRolePermissionWorkflow tests the complete role and permission workflow
func TestRolePermissionWorkflow(t *testing.T) {
	router, _, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Create organization
	t.Run("Create Organization", func(t *testing.T) {
		orgRequest := handlers.CreateOrganizationRequest{
			Name: "Test Organization",
		}

		body, _ := json.Marshal(orgRequest)
		req := httptest.NewRequest("POST", "/api/v1/organizations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response handlers.CreateOrganizationResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
	})

	// Test different permission combinations
	permissionTests := []struct {
		name        string
		permissions []models.Permission
		description string
	}{
		{
			name:        "read-only",
			permissions: []models.Permission{models.Read},
			description: "Can only read data",
		},
		{
			name:        "read-write",
			permissions: []models.Permission{models.Read, models.Write},
			description: "Can read and write data",
		},
		{
			name:        "full-access",
			permissions: []models.Permission{models.Read, models.Write, models.Delete},
			description: "Full access to all operations",
		},
		{
			name:        "write-only",
			permissions: []models.Permission{models.Write},
			description: "Can only write data (unusual but valid)",
		},
		{
			name:        "delete-only",
			permissions: []models.Permission{models.Delete},
			description: "Can only delete data (unusual but valid)",
		},
	}

	for _, tt := range permissionTests {
		t.Run("Create Role with "+tt.name+" permissions", func(t *testing.T) {
			roleRequest := handlers.CreateRoleRequest{
				Name:        tt.name,
				Permissions: tt.permissions,
			}

			body, _ := json.Marshal(roleRequest)
			req := httptest.NewRequest("POST", "/api/v1/roles", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusCreated, w.Code)

			var response handlers.CreateRoleResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tt.name, response.Role.Name)
			assert.Equal(t, len(tt.permissions), len(response.Role.Permissions))

			// Verify all expected permissions are present
			for _, expectedPerm := range tt.permissions {
				found := false
				for _, actualPerm := range response.Role.Permissions {
					if actualPerm == expectedPerm {
						found = true
						break
					}
				}
				assert.True(t, found, "Permission %d should be present", expectedPerm)
			}
		})
	}
}
