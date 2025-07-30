// @title           Unicorn API
// @version         1.0
// @description     A comprehensive RESTful API for Unicorn services providing Identity and Access Management (IAM) functionality. This API supports user authentication, role-based access control, organization management, and JWT token handling.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Unicorn API Support
// @contact.url    https://github.com/your-org/unicorn-api
// @contact.email  support@unicorn-api.com

// @license.name  MIT License
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token. Example: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

// @tag.name auth
// @tag.description Authentication and authorization endpoints

// @tag.name users
// @tag.description User management endpoints

// @tag.name roles
// @tag.description Role and permission management endpoints

// @tag.name organizations
// @tag.description Organization management endpoints

// @tag.name health
// @tag.description Health check and monitoring endpoints

// @tag.name hello
// @tag.description Basic connectivity test endpoints
package main

import (
	"log"
	"os"

	"unicorn-api/internal/config"
	"unicorn-api/internal/handlers"
	"unicorn-api/internal/middleware"
	"unicorn-api/internal/routes"
	"unicorn-api/internal/stores"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "unicorn-api/docs" // Import the generated docs package

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
)

func setupServices(cfg *config.Config) *handlers.IAMHandler {
	// setup the stores
	store, err := stores.NewGORMIAMStore("test.db")
	if err != nil {
		panic("failed to initialize IAM store: " + err.Error())
	}

	iamHandler := handlers.NewIAMHandler(store, cfg)

	return iamHandler
}

func main() {
	log.Println("Starting Unicorn API...")

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize configuration
	cfg := config.New()

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()

	// Add middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	// setup services
	iamHandler := setupServices(cfg)

	// Setup routes
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	routes.SetupRoutes(router, iamHandler)
	router.GET("/health", handlers.HealthCheck)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Unicorn API server on port %s", port)
	log.Printf("Version: %s, Build Time: %s", Version, BuildTime)
	log.Printf("Environment: %s", cfg.Environment)
	log.Printf("Swagger UI available at: http://localhost:%s/swagger/index.html", port)

	// Start server
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
