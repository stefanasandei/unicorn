package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheckResponse represents the health check response
// swagger:model
type HealthCheckResponse struct {
	// The health status of the API
	// example: healthy
	Status string `json:"status"`
	// The timestamp when the health check was performed
	// example: 2024-01-01T12:00:00Z
	Timestamp time.Time `json:"timestamp"`
	// The version of the API
	// example: 1.0.0
	Version string `json:"version"`
	// Additional health information
	// example: {"database":"connected","redis":"connected"}
	Details map[string]interface{} `json:"details,omitempty"`
}

// HelloResponse represents the hello endpoint response
// swagger:model
type HelloResponse struct {
	// The greeting message
	// example: Hello, World!
	Message string `json:"message"`
	// The timestamp when the request was processed
	// example: 2024-01-01T12:00:00Z
	Timestamp time.Time `json:"timestamp"`
	// The API version
	// example: 1.0.0
	Version string `json:"version"`
}

// HealthCheck godoc
// @Summary      Health check
// @Description  Get the health status of the API. This endpoint can be used by load balancers and monitoring systems to verify the API is running properly.
// @Tags         health
// @Produce      json
// @Success      200   {object}  HealthCheckResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /health [get]
func HealthCheck(c *gin.Context) {
	response := HealthCheckResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
	}

	c.JSON(http.StatusOK, response)
}

// Hello godoc
// @Summary      Hello endpoint
// @Description  Returns a hello world message. This endpoint is useful for testing basic connectivity to the API.
// @Tags         hello
// @Produce      json
// @Success      200   {object}  HelloResponse
// @Router       /api/v1/hello [get]
func Hello(c *gin.Context) {
	response := HelloResponse{
		Message:   "Hello, World!",
		Timestamp: time.Now(),
		Version:   "1.0.0",
	}

	c.JSON(http.StatusOK, response)
}
