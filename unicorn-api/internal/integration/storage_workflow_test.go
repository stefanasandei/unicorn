package integration

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"unicorn-api/internal/config"
	"unicorn-api/internal/handlers"
	"unicorn-api/internal/models"
	"unicorn-api/internal/routes"
	"unicorn-api/internal/stores"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupStorageIntegrationTest(t *testing.T) (*gin.Engine, func()) {
	// Create temporary database files
	iamDBPath := "storage_integration_iam_test_" + uuid.New().String() + ".db"
	storageDBPath := "storage_integration_storage_test_" + uuid.New().String() + ".db"
	storagePath := "storage_test_" + uuid.New().String()

	// Create storage directory
	err := os.MkdirAll(storagePath, 0755)
	require.NoError(t, err)

	// Setup IAM store
	iamStore, err := stores.NewGORMIAMStore(iamDBPath)
	require.NoError(t, err)

	// Setup storage store
	storageStore, err := stores.NewGORMStorageStore(storageDBPath, storagePath)
	require.NoError(t, err)

	cfg := &config.Config{
		JWTSecret:       "integration-test-secret-key",
		TokenExpiration: 24 * time.Hour,
		Environment:     "test",
	}

	iamHandler := handlers.NewIAMHandler(iamStore, cfg)
	storageHandler := handlers.NewStorageHandler(storageStore, iamStore, cfg)
	computeHandler := handlers.NewComputeHandler(cfg, iamStore)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	routes.SetupRoutes(router, iamHandler, storageHandler, computeHandler)

	// Return cleanup function
	cleanup := func() {
		// Drop tables and cleanup
		os.Remove(iamDBPath)
		os.Remove(storageDBPath)
		os.RemoveAll(storagePath)
	}

	return router, cleanup
}

// TestFullStorageWorkflow tests the complete storage workflow:
// 1. Create organization and roles
// 2. Create users with different permissions
// 3. Test bucket creation with different permission levels
// 4. Test file operations with proper authorization
func TestFullStorageWorkflow(t *testing.T) {
	router, cleanup := setupStorageIntegrationTest(t)
	defer cleanup()

	// Step 1: Create organization
	var orgID string
	t.Run("Create Organization", func(t *testing.T) {
		orgRequest := handlers.CreateOrganizationRequest{
			Name: "Storage Test Corp",
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

		orgID = response.Organization.ID.String()
	})

	// Step 2: Create roles with different permissions
	var adminRoleID, userRoleID, readonlyRoleID string
	t.Run("Create Admin Role", func(t *testing.T) {
		roleRequest := handlers.CreateRoleRequest{
			Name:        "storage_admin",
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

		adminRoleID = response.Role.ID.String()
	})

	t.Run("Create User Role", func(t *testing.T) {
		roleRequest := handlers.CreateRoleRequest{
			Name:        "storage_user",
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

		userRoleID = response.Role.ID.String()
	})

	t.Run("Create Read-Only Role", func(t *testing.T) {
		roleRequest := handlers.CreateRoleRequest{
			Name:        "storage_readonly",
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

		readonlyRoleID = response.Role.ID.String()
	})

	// Step 3: Create users with different roles
	var adminUser, regularUser, readonlyUser *loginResponse
	t.Run("Create Admin User", func(t *testing.T) {
		userRequest := handlers.CreateUserRequest{
			Name:     "Storage Admin",
			Email:    "admin@storage.com",
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

		// Login to get token
		loginReq := handlers.LoginRequest{
			Email:    "admin@storage.com",
			Password: "adminpass123",
		}
		adminUser = loginUser(t, router, loginReq)
	})

	t.Run("Create Regular User", func(t *testing.T) {
		userRequest := handlers.CreateUserRequest{
			Name:     "Storage User",
			Email:    "user@storage.com",
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

		// Login to get token
		loginReq := handlers.LoginRequest{
			Email:    "user@storage.com",
			Password: "userpass123",
		}
		regularUser = loginUser(t, router, loginReq)
	})

	t.Run("Create Read-Only User", func(t *testing.T) {
		userRequest := handlers.CreateUserRequest{
			Name:     "Storage ReadOnly",
			Email:    "readonly@storage.com",
			Password: "readonly123",
			RoleID:   readonlyRoleID,
		}

		body, _ := json.Marshal(userRequest)
		url := "/api/v1/organizations/" + orgID + "/users"
		req := httptest.NewRequest("POST", url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Login to get token
		loginReq := handlers.LoginRequest{
			Email:    "readonly@storage.com",
			Password: "readonly123",
		}
		readonlyUser = loginUser(t, router, loginReq)
	})

	// Step 4: Test bucket creation with different permissions
	var adminBucketID, userBucketID string
	t.Run("Admin creates bucket", func(t *testing.T) {
		bucketRequest := handlers.CreateBucketRequest{
			Name: "admin-bucket",
		}

		body, _ := json.Marshal(bucketRequest)
		req := httptest.NewRequest("POST", "/api/v1/buckets", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+adminUser.Token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var bucket models.StorageBucket
		err := json.Unmarshal(w.Body.Bytes(), &bucket)
		require.NoError(t, err)
		assert.Equal(t, "admin-bucket", bucket.Name)
		adminBucketID = bucket.ID.String()
	})

	t.Run("User creates bucket", func(t *testing.T) {
		bucketRequest := handlers.CreateBucketRequest{
			Name: "user-bucket",
		}

		body, _ := json.Marshal(bucketRequest)
		req := httptest.NewRequest("POST", "/api/v1/buckets", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+regularUser.Token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var bucket models.StorageBucket
		err := json.Unmarshal(w.Body.Bytes(), &bucket)
		require.NoError(t, err)
		assert.Equal(t, "user-bucket", bucket.Name)
		userBucketID = bucket.ID.String()
	})

	t.Run("ReadOnly user cannot create bucket", func(t *testing.T) {
		bucketRequest := handlers.CreateBucketRequest{
			Name: "readonly-bucket",
		}

		body, _ := json.Marshal(bucketRequest)
		req := httptest.NewRequest("POST", "/api/v1/buckets", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+readonlyUser.Token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	// Step 5: Test file operations
	t.Run("Upload file to bucket", func(t *testing.T) {
		fileContent := "Hello, World!"
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "test.txt")
		require.NoError(t, err)
		_, err = part.Write([]byte(fileContent))
		require.NoError(t, err)
		writer.Close()

		url := "/api/v1/buckets/" + adminBucketID + "/files"
		req := httptest.NewRequest("POST", url, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+adminUser.Token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var file models.File
		err = json.Unmarshal(w.Body.Bytes(), &file)
		require.NoError(t, err)
		assert.Equal(t, "test.txt", file.Name)
		assert.Equal(t, int64(13), file.Size)
	})

	t.Run("List files in bucket", func(t *testing.T) {
		url := "/api/v1/buckets/" + adminBucketID + "/files"
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+adminUser.Token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var files []models.File
		err := json.Unmarshal(w.Body.Bytes(), &files)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(files), 1)
	})

	t.Run("Download file from bucket", func(t *testing.T) {
		// First upload a file
		fileContent := "Download test content"
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "download.txt")
		require.NoError(t, err)
		_, err = part.Write([]byte(fileContent))
		require.NoError(t, err)
		writer.Close()

		uploadURL := "/api/v1/buckets/" + userBucketID + "/files"
		uploadReq := httptest.NewRequest("POST", uploadURL, body)
		uploadReq.Header.Set("Content-Type", writer.FormDataContentType())
		uploadReq.Header.Set("Authorization", "Bearer "+regularUser.Token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, uploadReq)
		assert.Equal(t, http.StatusCreated, w.Code)

		var uploadedFile models.File
		err = json.Unmarshal(w.Body.Bytes(), &uploadedFile)
		require.NoError(t, err)

		// Now download the file
		downloadURL := "/api/v1/buckets/" + userBucketID + "/files/" + uploadedFile.ID.String()
		downloadReq := httptest.NewRequest("GET", downloadURL, nil)
		downloadReq.Header.Set("Authorization", "Bearer "+regularUser.Token)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, downloadReq)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, fileContent, w.Body.String())
	})

	t.Run("User cannot access another user's bucket", func(t *testing.T) {
		url := "/api/v1/buckets/" + adminBucketID + "/files"
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+regularUser.Token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Delete file from bucket", func(t *testing.T) {
		// First upload a file
		fileContent := "File to delete"
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "delete.txt")
		require.NoError(t, err)
		_, err = part.Write([]byte(fileContent))
		require.NoError(t, err)
		writer.Close()

		uploadURL := "/api/v1/buckets/" + adminBucketID + "/files"
		uploadReq := httptest.NewRequest("POST", uploadURL, body)
		uploadReq.Header.Set("Content-Type", writer.FormDataContentType())
		uploadReq.Header.Set("Authorization", "Bearer "+adminUser.Token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, uploadReq)
		assert.Equal(t, http.StatusCreated, w.Code)

		var uploadedFile models.File
		err = json.Unmarshal(w.Body.Bytes(), &uploadedFile)
		require.NoError(t, err)

		// Delete the file
		deleteURL := "/api/v1/buckets/" + adminBucketID + "/files/" + uploadedFile.ID.String()
		deleteReq := httptest.NewRequest("DELETE", deleteURL, nil)
		deleteReq.Header.Set("Authorization", "Bearer "+adminUser.Token)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, deleteReq)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("ReadOnly user cannot delete file", func(t *testing.T) {
		// First upload a file as admin
		fileContent := "File readonly cannot delete"
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "readonly_delete.txt")
		require.NoError(t, err)
		_, err = part.Write([]byte(fileContent))
		require.NoError(t, err)
		writer.Close()

		uploadURL := "/api/v1/buckets/" + adminBucketID + "/files"
		uploadReq := httptest.NewRequest("POST", uploadURL, body)
		uploadReq.Header.Set("Content-Type", writer.FormDataContentType())
		uploadReq.Header.Set("Authorization", "Bearer "+adminUser.Token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, uploadReq)
		assert.Equal(t, http.StatusCreated, w.Code)

		var uploadedFile models.File
		err = json.Unmarshal(w.Body.Bytes(), &uploadedFile)
		require.NoError(t, err)

		// Try to delete as readonly user
		deleteURL := "/api/v1/buckets/" + adminBucketID + "/files/" + uploadedFile.ID.String()
		deleteReq := httptest.NewRequest("DELETE", deleteURL, nil)
		deleteReq.Header.Set("Authorization", "Bearer "+readonlyUser.Token)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, deleteReq)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

// loginResponse represents the login response

type loginResponse struct {
	Token string `json:"token"`
}

// loginUser helper function to login and get token
func loginUser(t *testing.T, router *gin.Engine, loginReq handlers.LoginRequest) *loginResponse {
	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response loginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	return &response
}

// TestStorageAuthorization tests authorization scenarios
func TestStorageAuthorization(t *testing.T) {
	router, cleanup := setupStorageIntegrationTest(t)
	defer cleanup()

	// Setup test data
	orgID := createTestOrganization(t, router, "Auth Test Org")
	adminRoleID := createTestRole(t, router, "auth_admin", []models.Permission{models.Read, models.Write, models.Delete})
	userRoleID := createTestRole(t, router, "auth_user", []models.Permission{models.Read, models.Write})

	adminUser := createTestUser(t, router, orgID, adminRoleID, "authadmin@example.com", "adminpass")
	_ = createTestUser(t, router, orgID, userRoleID, "authuser@example.com", "userpass")

	// Test unauthorized access
	t.Run("Unauthorized access to buckets", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/buckets", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// Test invalid token
	t.Run("Invalid token access", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/buckets", nil)
		req.Header.Set("Authorization", "Bearer invalid_token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// Test expired token handling (would need to implement token expiration)
	t.Run("Valid token access", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/buckets", nil)
		req.Header.Set("Authorization", "Bearer "+adminUser.Token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// Helper functions for test setup
func createTestOrganization(t *testing.T, router *gin.Engine, name string) string {
	orgRequest := handlers.CreateOrganizationRequest{Name: name}
	body, _ := json.Marshal(orgRequest)
	req := httptest.NewRequest("POST", "/api/v1/organizations", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response handlers.CreateOrganizationResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	return response.Organization.ID.String()
}

func createTestRole(t *testing.T, router *gin.Engine, name string, permissions []models.Permission) string {
	roleRequest := handlers.CreateRoleRequest{Name: name, Permissions: permissions}
	body, _ := json.Marshal(roleRequest)
	req := httptest.NewRequest("POST", "/api/v1/roles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response handlers.CreateRoleResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	return response.Role.ID.String()
}

func createTestUser(t *testing.T, router *gin.Engine, orgID, roleID, email, password string) *loginResponse {
	userRequest := handlers.CreateUserRequest{
		Name:     "Test User",
		Email:    email,
		Password: password,
		RoleID:   roleID,
	}

	body, _ := json.Marshal(userRequest)
	url := "/api/v1/organizations/" + orgID + "/users"
	req := httptest.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	loginReq := handlers.LoginRequest{
		Email:    email,
		Password: password,
	}
	return loginUser(t, router, loginReq)
}
