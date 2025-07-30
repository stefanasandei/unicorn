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
// @type apiKey
// @description "Type 'Bearer ' + your JWT token to authorize"
package main

import (
	"log"
	"os"
	"path/filepath"

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

// setupServices initializes all services and handlers
func setupServices(cfg *config.Config) (*handlers.IAMHandler, *handlers.StorageHandler, *handlers.ComputeHandler, *handlers.LambdaHandler, *handlers.SecretsHandler) {
	// Get database path from environment or use default
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "unicorn.db"
	}

	// Ensure database directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatal("Failed to create database directory:", err)
	}

	// Setup IAM store
	iamStore, err := stores.NewGORMIAMStore(dbPath)
	if err != nil {
		log.Fatal("Failed to initialize IAM store:", err)
	}
	if err := iamStore.SeedAdmin(cfg); err != nil {
		log.Fatal("Failed to seed admin:", err)
	}

	// Setup storage store
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./storage"
	}
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		log.Fatal("Failed to create storage directory:", err)
	}
	storageStore, err := stores.NewGORMStorageStore(dbPath, storagePath)
	if err != nil {
		log.Fatal("Failed to initialize storage store:", err)
	}

	// Setup secrets store
	secretsStore, err := stores.NewSecretStore(dbPath)
	if err != nil {
		log.Fatal("Failed to initialize secrets store:", err)
	}

	// Create handlers
	iamHandler := handlers.NewIAMHandler(iamStore, cfg)
	secretsHandler := handlers.NewSecretsHandler(secretsStore, iamStore, cfg)
	storageHandler := handlers.NewStorageHandler(storageStore, iamStore, cfg)
	computeHandler := handlers.NewComputeHandler(cfg, iamStore)
	lambdaHandler := handlers.NewLambdaHandler(cfg, iamStore)

	return iamHandler, storageHandler, computeHandler, lambdaHandler, secretsHandler
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

	// Setup services
	iamHandler, storageHandler, computeHandler, lambdaHandler, secretsHandler := setupServices(cfg)

	// Setup routes
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	routes.SetupRoutes(router, iamHandler, storageHandler, computeHandler, lambdaHandler, secretsHandler, cfg)
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
