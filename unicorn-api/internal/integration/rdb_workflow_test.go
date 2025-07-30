package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"unicorn-api/internal/config"
	"unicorn-api/internal/handlers"
	"unicorn-api/internal/models"
	"unicorn-api/internal/routes"
	"unicorn-api/internal/stores"
)

func setupRDBTestServer() *gin.Engine {
	_ = os.Remove("test.db") // Ensure fresh DB for each test
	cfg := config.New()
	cfg.JWTSecret = "test-secret" // Ensure consistent JWT secret for test

	store, _ := stores.NewGORMIAMStore("test.db")
	_ = store.SeedAdmin(cfg)
	iamHandler := handlers.NewIAMHandler(store, cfg)
	computeHandler := handlers.NewComputeHandler(cfg, store)
	secretsStore, err := stores.NewSecretStore("test.db")
	if err != nil {
		panic("failed to initialize secrets store: " + err.Error())
	}
	storageHandler := handlers.NewStorageHandler(&stores.GORMStorageStore{}, store, cfg)
	lambdaHandler := handlers.NewLambdaHandler(cfg, store)
	secretsHandler := handlers.NewSecretsHandler(secretsStore, store, cfg)
	rdbHandler := handlers.NewRDBHandler(cfg, store)

	router := gin.Default()
	routes.SetupRoutes(router, iamHandler, storageHandler, computeHandler, lambdaHandler, secretsHandler, rdbHandler, cfg)
	return router
}

func TestRDBCreateAndList(t *testing.T) {
	if os.Getenv("SKIP_DOCKER_TESTS") == "1" {
		t.Skip("Skipping Docker integration test")
	}

	// Check if Docker is available by trying to create a client
	// If Docker is not available, skip the test
	if os.Getenv("DOCKER_HOST") == "" {
		// Try to detect if Docker is running
		_, err := os.Stat("/var/run/docker.sock")
		if err != nil {
			t.Skip("Docker not available - skipping RDB test")
		}
	}
	router := setupRDBTestServer()

	// Login as admin to get JWT token
	loginReq := map[string]string{"email": "admin@unicorn.local", "password": "admin123"}
	loginBody, _ := json.Marshal(loginReq)
	login := httptest.NewRequest("POST", "/api/v1/login", bytes.NewReader(loginBody))
	login.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, login)
	assert.Equal(t, http.StatusOK, loginW.Code)
	var loginResp map[string]interface{}
	_ = json.Unmarshal(loginW.Body.Bytes(), &loginResp)
	token, ok := loginResp["token"].(string)
	assert.True(t, ok)

	// Test PostgreSQL instance creation
	t.Run("Create PostgreSQL Instance", func(t *testing.T) {
		createReq := models.RDBCreateRequest{
			Name:   "test-postgres",
			Type:   models.RDBTypePostgreSQL,
			Preset: models.RDBPresetMicro,
			Database: "testdb",
			Username: "testuser",
			Password: "testpass123",
			Volumes: []models.RDBVolume{
				{
					Name:      "data",
					Size:      1,
					MountPath: "/var/lib/postgresql/data",
				},
			},
		}
		body, _ := json.Marshal(createReq)

		req := httptest.NewRequest("POST", "/api/v1/rdb/create", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var instanceInfo models.RDBInstanceInfo
		_ = json.Unmarshal(w.Body.Bytes(), &instanceInfo)
		assert.Equal(t, "test-postgres", instanceInfo.Name)
		assert.Equal(t, models.RDBTypePostgreSQL, instanceInfo.Type)
		assert.Equal(t, "testdb", instanceInfo.Database)
		assert.Equal(t, "testuser", instanceInfo.Username)
		assert.Len(t, instanceInfo.Volumes, 1)
		assert.Equal(t, "data", instanceInfo.Volumes[0].Name)
	})

	// Test MySQL instance creation
	t.Run("Create MySQL Instance", func(t *testing.T) {
		createReq := models.RDBCreateRequest{
			Name:   "test-mysql",
			Type:   models.RDBTypeMySQL,
			Preset: models.RDBPresetSmall,
			Database: "testdb",
			Username: "testuser",
			Password: "testpass123",
		}
		body, _ := json.Marshal(createReq)

		req := httptest.NewRequest("POST", "/api/v1/rdb/create", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var instanceInfo models.RDBInstanceInfo
		_ = json.Unmarshal(w.Body.Bytes(), &instanceInfo)
		assert.Equal(t, "test-mysql", instanceInfo.Name)
		assert.Equal(t, models.RDBTypeMySQL, instanceInfo.Type)
		assert.Equal(t, "testdb", instanceInfo.Database)
		assert.Equal(t, "testuser", instanceInfo.Username)
	})

	// Test listing RDB instances
	t.Run("List RDB Instances", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/rdb/list", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var instances []models.RDBInstanceInfo
		_ = json.Unmarshal(w.Body.Bytes(), &instances)
		assert.GreaterOrEqual(t, len(instances), 2) // Should have at least the 2 instances we created
	})
}

func TestRDBPermissionChecks(t *testing.T) {
	if os.Getenv("SKIP_DOCKER_TESTS") == "1" {
		t.Skip("Skipping Docker integration test")
	}

	router := setupRDBTestServer()

	// Create organization and roles for permission testing
	var orgID, readonlyRoleID string

	// Create organization
	orgRequest := handlers.CreateOrganizationRequest{Name: "RDB Test Corp"}
	orgBody, _ := json.Marshal(orgRequest)
	orgReq := httptest.NewRequest("POST", "/api/v1/organizations", bytes.NewReader(orgBody))
	orgReq.Header.Set("Content-Type", "application/json")
	orgW := httptest.NewRecorder()
	router.ServeHTTP(orgW, orgReq)
	assert.Equal(t, http.StatusCreated, orgW.Code)

	var orgResp handlers.CreateOrganizationResponse
	_ = json.Unmarshal(orgW.Body.Bytes(), &orgResp)
	orgID = orgResp.Organization.ID.String()

	// Create readonly role
	readonlyRoleReq := handlers.CreateRoleRequest{
		Name:        "rdb_readonly",
		Permissions: []models.Permission{models.Read},
	}
	readonlyRoleBody, _ := json.Marshal(readonlyRoleReq)
	readonlyRoleRequest := httptest.NewRequest("POST", "/api/v1/roles", bytes.NewReader(readonlyRoleBody))
	readonlyRoleRequest.Header.Set("Content-Type", "application/json")
	readonlyRoleW := httptest.NewRecorder()
	router.ServeHTTP(readonlyRoleW, readonlyRoleRequest)
	assert.Equal(t, http.StatusCreated, readonlyRoleW.Code)

	var readonlyRoleResp handlers.CreateRoleResponse
	_ = json.Unmarshal(readonlyRoleW.Body.Bytes(), &readonlyRoleResp)
	readonlyRoleID = readonlyRoleResp.Role.ID.String()

	// Create write role (not used in this test but kept for completeness)
	writeRoleReq := handlers.CreateRoleRequest{
		Name:        "rdb_write",
		Permissions: []models.Permission{models.Read, models.Write},
	}
	writeRoleBody, _ := json.Marshal(writeRoleReq)
	writeRoleRequest := httptest.NewRequest("POST", "/api/v1/roles", bytes.NewReader(writeRoleBody))
	writeRoleRequest.Header.Set("Content-Type", "application/json")
	writeRoleW := httptest.NewRecorder()
	router.ServeHTTP(writeRoleW, writeRoleRequest)
	assert.Equal(t, http.StatusCreated, writeRoleW.Code)

	// Create readonly user
	readonlyUserReq := handlers.CreateUserRequest{
		RoleID:   readonlyRoleID,
		Name:     "Readonly User",
		Email:    "readonly@test.com",
		Password: "password123",
	}
	readonlyUserBody, _ := json.Marshal(readonlyUserReq)
	readonlyUserRequest := httptest.NewRequest("POST", "/api/v1/organizations/"+orgID+"/users", bytes.NewReader(readonlyUserBody))
	readonlyUserRequest.Header.Set("Content-Type", "application/json")
	readonlyUserW := httptest.NewRecorder()
	router.ServeHTTP(readonlyUserW, readonlyUserRequest)
	assert.Equal(t, http.StatusCreated, readonlyUserW.Code)

	// Login as readonly user
	readonlyLoginReq := map[string]string{"email": "readonly@test.com", "password": "password123"}
	readonlyLoginBody, _ := json.Marshal(readonlyLoginReq)
	readonlyLogin := httptest.NewRequest("POST", "/api/v1/login", bytes.NewReader(readonlyLoginBody))
	readonlyLogin.Header.Set("Content-Type", "application/json")
	readonlyLoginW := httptest.NewRecorder()
	router.ServeHTTP(readonlyLoginW, readonlyLogin)
	assert.Equal(t, http.StatusOK, readonlyLoginW.Code)

	var readonlyLoginResp map[string]interface{}
	_ = json.Unmarshal(readonlyLoginW.Body.Bytes(), &readonlyLoginResp)
	readonlyToken, _ := readonlyLoginResp["token"].(string)

	// Test that readonly user can list but not create
	t.Run("Readonly User Can List But Not Create", func(t *testing.T) {
		// Should be able to list
		listReq := httptest.NewRequest("GET", "/api/v1/rdb/list", nil)
		listReq.Header.Set("Authorization", "Bearer "+readonlyToken)
		listW := httptest.NewRecorder()
		router.ServeHTTP(listW, listReq)
		assert.Equal(t, http.StatusOK, listW.Code)

		// Should not be able to create
		createReq := models.RDBCreateRequest{
			Name:   "test-instance",
			Type:   models.RDBTypePostgreSQL,
			Preset: models.RDBPresetMicro,
		}
		createBody, _ := json.Marshal(createReq)
		createRequest := httptest.NewRequest("POST", "/api/v1/rdb/create", bytes.NewReader(createBody))
		createRequest.Header.Set("Content-Type", "application/json")
		createRequest.Header.Set("Authorization", "Bearer "+readonlyToken)
		createW := httptest.NewRecorder()
		router.ServeHTTP(createW, createRequest)
		assert.Equal(t, http.StatusForbidden, createW.Code)
	})
} 