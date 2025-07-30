package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"unicorn-api/internal/auth"
	"unicorn-api/internal/config"
	"unicorn-api/internal/handlers"
	"unicorn-api/internal/middleware"
	"unicorn-api/internal/models"
	"unicorn-api/internal/stores"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func setupTestSecretsHandler(t *testing.T) (*handlers.SecretsHandler, *gin.Engine, uuid.UUID, string) {
	cfg := config.New()
	secretStore, err := stores.NewSecretStore(":memory:")
	require.NoError(t, err)

	// Create a mock IAM store
	iamStore, err := stores.NewGORMIAMStore(":memory:")
	require.NoError(t, err)

	// Create an admin role with all permissions
	role := &models.Role{
		ID:          uuid.New(),
		Name:        "admin",
		Permissions: models.Permissions{models.Read, models.Write, models.Delete},
	}
	err = iamStore.CreateRole(role)
	require.NoError(t, err)

	// Create a test user with the admin role
	userID := uuid.New()
	token, err := auth.GenerateToken(userID.String(), role.ID.String(), cfg)
	require.NoError(t, err)

	h := &handlers.SecretsHandler{
		Store:    secretStore,
		Config:   cfg,
		IAMStore: iamStore,
	}

	// Setup router with proper authentication middleware
	router := gin.Default()

	// Add authentication middleware to protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.GET("/secrets", h.ListSecrets)
		protected.POST("/secrets", h.CreateSecret)
		protected.GET("/secrets/:id", h.ReadSecret)
		protected.PUT("/secrets/:id", h.UpdateSecret)
		protected.DELETE("/secrets/:id", h.DeleteSecret)
	}

	return h, router, userID, token
}

func TestSecretsCRUD(t *testing.T) {
	_, r, _, token := setupTestSecretsHandler(t)

	// Create
	body := map[string]interface{}{
		"name":     "api-key",
		"value":    "supersecret",
		"metadata": "{\"env\":\"test\"}",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/secrets", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	secretID := resp["id"].(string)

	// List
	req = httptest.NewRequest("GET", "/api/v1/secrets", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	var list []map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &list))
	require.Len(t, list, 1)

	// Read
	url := "/api/v1/secrets/" + secretID
	req = httptest.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	var read map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &read))
	require.Equal(t, "supersecret", read["value"])

	// Update
	update := map[string]interface{}{"value": "newsecret"}
	b, _ = json.Marshal(update)
	req = httptest.NewRequest("PUT", url, bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)

	// Read updated
	req = httptest.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &read))
	require.Equal(t, "newsecret", read["value"])

	// Delete
	req = httptest.NewRequest("DELETE", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)

	// List after delete
	req = httptest.NewRequest("GET", "/api/v1/secrets", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	list = nil
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &list))
	require.Len(t, list, 0)
}
