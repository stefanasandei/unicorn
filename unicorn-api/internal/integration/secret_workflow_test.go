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
	"unicorn-api/internal/stores"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func setupTestSecretsHandler(t *testing.T) (*handlers.SecretsHandler, *gin.Engine, uuid.UUID, string) {
	cfg := config.New()
	store, err := stores.NewSecretStore(":memory:")
	require.NoError(t, err)
	userID := uuid.New()
	token, err := auth.GenerateToken(userID.String(), "role", cfg)
	require.NoError(t, err)
	h := &handlers.SecretsHandler{Store: store, Config: cfg}
	r := gin.Default()
	r.GET("/secrets", h.ListSecrets)
	r.POST("/secrets", h.CreateSecret)
	r.GET("/secrets/:id", h.ReadSecret)
	r.PUT("/secrets/:id", h.UpdateSecret)
	r.DELETE("/secrets/:id", h.DeleteSecret)
	return h, r, userID, token
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
	req := httptest.NewRequest("POST", "/secrets", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	secretID := resp["id"].(string)

	// List
	req = httptest.NewRequest("GET", "/secrets", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	var list []map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &list))
	require.Len(t, list, 1)

	// Read
	url := "/secrets/" + secretID
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
	req = httptest.NewRequest("GET", "/secrets", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	list = nil
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &list))
	require.Len(t, list, 0)
}
