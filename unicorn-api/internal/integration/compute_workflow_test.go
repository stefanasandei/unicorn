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

func setupTestServer() *gin.Engine {
	_ = os.Remove("test.db") // Ensure fresh DB for each test
	cfg := config.New()
	cfg.JWTSecret = "test-secret" // Ensure consistent JWT secret for test
	store, _ := stores.NewGORMIAMStore("test.db")
	secretsStore, err := stores.NewSecretStore("test.db")
	if err != nil {
		panic("failed to initialize secrets store: " + err.Error())
	}

	_ = store.SeedAdmin(cfg)
	iamHandler := handlers.NewIAMHandler(store, cfg)
	secretsHandler := handlers.NewSecretsHandler(secretsStore, store, cfg)
	computeHandler := handlers.NewComputeHandler(cfg, store)
	storageHandler := handlers.NewStorageHandler(&stores.GORMStorageStore{}, store, cfg)
	lambdaHandler := handlers.NewLambdaHandler(cfg, store)
	rdbHandler := handlers.NewRDBHandler(cfg, store)
	router := gin.Default()
	routes.SetupRoutes(router, iamHandler, storageHandler, computeHandler, lambdaHandler, secretsHandler, rdbHandler, cfg)
	return router
}

func TestComputeCreateAndList(t *testing.T) {
	if os.Getenv("SKIP_DOCKER_TESTS") == "1" {
		t.Skip("Skipping Docker integration test")
	}

	// Check if Docker is available by trying to create a client
	// If Docker is not available, skip the test
	if os.Getenv("DOCKER_HOST") == "" {
		// Try to detect if Docker is running
		_, err := os.Stat("/var/run/docker.sock")
		if err != nil {
			t.Skip("Docker not available - skipping compute test")
		}
	}
	router := setupTestServer()

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

	// Prepare request
	createReq := models.ComputeCreateRequest{
		Image:      "nginx:alpine",
		Preset:     models.PresetMicro,
		Ports:      map[string]string{"80": "8081"},
		ExposePort: "80",
	}
	body, _ := json.Marshal(createReq)

	req := httptest.NewRequest("POST", "/api/v1/compute/create", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// List containers
	req2 := httptest.NewRequest("GET", "/api/v1/compute/list", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	var containers []models.ComputeContainerInfo
	_ = json.Unmarshal(w2.Body.Bytes(), &containers)
	assert.GreaterOrEqual(t, len(containers), 1)
}
