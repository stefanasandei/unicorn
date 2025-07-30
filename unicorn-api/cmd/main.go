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

func setupServices(cfg *config.Config) (*handlers.IAMHandler, *handlers.StorageHandler, *handlers.ComputeHandler) {
	// setup the stores
	store, err := stores.NewGORMIAMStore("test.db")
	if err != nil {
		panic("failed to initialize IAM store: " + err.Error())
	}
	if err := store.SeedAdmin(); err != nil {
		panic("failed to seed admin: " + err.Error())
	}

	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./storage" // default fallback
	}
	storageStore, err := stores.NewGORMStorageStore("test.db", storagePath)
	if err != nil {
		panic("failed to initialize storage store: " + err.Error())
	}

	iamHandler := handlers.NewIAMHandler(store, cfg)
	storageHandler := handlers.NewStorageHandler(storageStore, store, cfg)
	computeHandler := handlers.NewComputeHandler(cfg, store)

	return iamHandler, storageHandler, computeHandler
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
	iamHandler, storageHandler, computeHandler := setupServices(cfg)

	// Setup routes
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	routes.SetupRoutes(router, iamHandler, storageHandler, computeHandler)
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
