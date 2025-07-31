package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"unicorn-api/internal/models"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Open SQLite database
	db, err := sql.Open("sqlite3", "./unicorn.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Create monitoring tables if they don't exist
	if err := createTables(db); err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	// Generate sample data
	if err := populateSampleData(db); err != nil {
		log.Fatal("Failed to populate sample data:", err)
	}

	fmt.Println("Successfully populated monitoring database with sample data!")
}

func createTables(db *sql.DB) error {
	// Create resource_usages table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS resource_usages (
			id TEXT PRIMARY KEY,
			account_id TEXT NOT NULL,
			organization_id TEXT NOT NULL,
			resource_type TEXT NOT NULL,
			resource_id TEXT NOT NULL,
			resource_name TEXT NOT NULL,
			status TEXT NOT NULL,
			cpu_usage REAL DEFAULT 0,
			memory_usage INTEGER DEFAULT 0,
			storage_usage INTEGER DEFAULT 0,
			network_usage INTEGER DEFAULT 0,
			request_count INTEGER DEFAULT 0,
			execution_time REAL DEFAULT 0,
			cost_per_hour REAL NOT NULL,
			total_cost REAL DEFAULT 0,
			currency TEXT DEFAULT 'USD',
			resource_created_at DATETIME NOT NULL,
			last_active_at DATETIME,
			last_updated_at DATETIME NOT NULL,
			configuration TEXT,
			tags TEXT,
			notes TEXT
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create resource_usages table: %w", err)
	}

	// Create resource_usage_histories table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS resource_usage_histories (
			id TEXT PRIMARY KEY,
			account_id TEXT NOT NULL,
			organization_id TEXT NOT NULL,
			resource_type TEXT NOT NULL,
			resource_id TEXT NOT NULL,
			resource_name TEXT NOT NULL,
			avg_cpu_usage REAL DEFAULT 0,
			peak_cpu_usage REAL DEFAULT 0,
			avg_memory_usage INTEGER DEFAULT 0,
			peak_memory_usage INTEGER DEFAULT 0,
			total_storage_usage INTEGER DEFAULT 0,
			total_network_usage INTEGER DEFAULT 0,
			total_requests INTEGER DEFAULT 0,
			total_execution_time REAL DEFAULT 0,
			total_cost REAL DEFAULT 0,
			currency TEXT DEFAULT 'USD',
			period_start DATETIME NOT NULL,
			period_end DATETIME NOT NULL,
			duration_hours REAL DEFAULT 0,
			created_at DATETIME NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create resource_usage_histories table: %w", err)
	}

	// Create billing_periods table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS billing_periods (
			id TEXT PRIMARY KEY,
			organization_id TEXT NOT NULL,
			period_start DATETIME NOT NULL,
			period_end DATETIME NOT NULL,
			total_cost REAL DEFAULT 0,
			currency TEXT DEFAULT 'USD',
			is_paid BOOLEAN DEFAULT FALSE,
			payment_method TEXT,
			transaction_id TEXT,
			compute_cost REAL DEFAULT 0,
			lambda_cost REAL DEFAULT 0,
			storage_cost REAL DEFAULT 0,
			rdb_cost REAL DEFAULT 0,
			secret_cost REAL DEFAULT 0,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create billing_periods table: %w", err)
	}

	// Create monitoring_metrics table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS monitoring_metrics (
			id TEXT PRIMARY KEY,
			organization_id TEXT NOT NULL,
			resource_id TEXT NOT NULL,
			resource_type TEXT NOT NULL,
			cpu_usage REAL DEFAULT 0,
			memory_usage INTEGER DEFAULT 0,
			storage_usage INTEGER DEFAULT 0,
			network_usage INTEGER DEFAULT 0,
			request_count INTEGER DEFAULT 0,
			execution_time REAL DEFAULT 0,
			created_at DATETIME NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create monitoring_metrics table: %w", err)
	}

	return nil
}

func populateSampleData(db *sql.DB) error {
	// Use existing organization and account IDs from the database
	var orgID, accountID string

	// Get the first organization ID from accounts table
	err := db.QueryRow("SELECT organization_id FROM accounts LIMIT 1").Scan(&orgID)
	if err != nil {
		// If no accounts exist, create new ones
		orgID = uuid.New().String()
		accountID = uuid.New().String()
	} else {
		// Get the first account ID for this organization
		err = db.QueryRow("SELECT id FROM accounts WHERE organization_id = ? LIMIT 1", orgID).Scan(&accountID)
		if err != nil {
			accountID = uuid.New().String()
		}
	}

	// Sample resource data with realistic usage patterns
	resources := []struct {
		resourceType models.ResourceType
		resourceName string
		costPerHour  float64
		baseUsage    int
	}{
		{models.ResourceTypeCompute, "web-server-1", 0.05, 45},
		{models.ResourceTypeCompute, "api-server-1", 0.10, 65},
		{models.ResourceTypeStorage, "data-bucket-1", 0.01, 0},
		{models.ResourceTypeStorage, "backup-bucket-1", 0.01, 0},
		{models.ResourceTypeLambda, "data-processor", 0.02, 0},
		{models.ResourceTypeLambda, "image-resizer", 0.02, 0},
		{models.ResourceTypeRDB, "main-database", 0.15, 0},
		{models.ResourceTypeSecret, "api-keys", 0.01, 0},
	}

	now := time.Now()
	rand.Seed(now.UnixNano())

	for i, resource := range resources {
		// Create resource usage
		resourceID := fmt.Sprintf("%s-%d", resource.resourceType, i+1)
		createdAt := now.Add(-time.Duration(rand.Intn(30)) * 24 * time.Hour) // Random creation time within last 30 days

		// Generate realistic usage metrics
		cpuUsage := float64(resource.baseUsage) + rand.Float64()*20 // Add some variance
		memoryUsage := rand.Intn(2048) + 512                        // 512MB to 2.5GB
		storageUsage := rand.Intn(10240) + 1024                     // 1GB to 11GB
		networkUsage := rand.Intn(512) + 128                        // 128MB to 640MB

		// Calculate total cost based on usage duration
		duration := now.Sub(createdAt)
		hours := duration.Hours()
		totalCost := resource.costPerHour * hours

		// Add storage and network costs
		storageCost := (float64(storageUsage) / 1024.0) * 0.01 * hours  // $0.01 per GB per hour
		networkCost := (float64(networkUsage) / 1024.0) * 0.001 * hours // $0.001 per GB per hour
		totalCost += storageCost + networkCost

		// Insert resource usage
		_, err := db.Exec(`
			INSERT OR REPLACE INTO resource_usages (
				id, account_id, organization_id, resource_type, resource_id, resource_name,
				status, cpu_usage, memory_usage, storage_usage, network_usage,
				cost_per_hour, total_cost, currency, resource_created_at, last_active_at, last_updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			uuid.New().String(), accountID, orgID, resource.resourceType, resourceID, resource.resourceName,
			models.ResourceStatusActive, cpuUsage, memoryUsage, storageUsage, networkUsage,
			resource.costPerHour, totalCost, "USD", createdAt, now, now,
		)
		if err != nil {
			return fmt.Errorf("failed to insert resource usage: %w", err)
		}

		// Create monitoring metrics
		_, err = db.Exec(`
			INSERT OR REPLACE INTO monitoring_metrics (
				id, organization_id, resource_id, resource_type, cpu_usage, memory_usage, storage_usage, network_usage, created_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			uuid.New().String(), orgID, resourceID, resource.resourceType, cpuUsage, memoryUsage, storageUsage, networkUsage, now,
		)
		if err != nil {
			return fmt.Errorf("failed to insert monitoring metrics: %w", err)
		}

		// Create historical usage records for the last 3 months
		for j := 0; j < 3; j++ {
			monthStart := time.Date(now.Year(), now.Month()-time.Month(j), 1, 0, 0, 0, 0, time.UTC)
			monthEnd := monthStart.AddDate(0, 1, -1)

			// Generate monthly usage data
			monthlyCPU := cpuUsage + rand.Float64()*10 - 5         // ±5% variance
			monthlyMemory := memoryUsage + rand.Intn(512) - 256    // ±256MB variance
			monthlyStorage := storageUsage + rand.Intn(1024) - 512 // ±512MB variance
			monthlyNetwork := networkUsage + rand.Intn(256) - 128  // ±128MB variance

			monthlyCost := resource.costPerHour * 24 * float64(monthEnd.Day()) // Full month cost
			monthlyCost += (float64(monthlyStorage) / 1024.0) * 0.01 * 24 * float64(monthEnd.Day())
			monthlyCost += (float64(monthlyNetwork) / 1024.0) * 0.001 * 24 * float64(monthEnd.Day())

			_, err := db.Exec(`
				INSERT OR REPLACE INTO resource_usage_histories (
					id, account_id, organization_id, resource_type, resource_id, resource_name,
					avg_cpu_usage, peak_cpu_usage, avg_memory_usage, peak_memory_usage,
					total_storage_usage, total_network_usage, total_cost, currency,
					period_start, period_end, duration_hours, created_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
				uuid.New().String(), accountID, orgID, resource.resourceType, resourceID, resource.resourceName,
				monthlyCPU, monthlyCPU*1.2, monthlyMemory, int(float64(monthlyMemory)*1.1),
				monthlyStorage, monthlyNetwork, monthlyCost, "USD",
				monthStart, monthEnd, 24*float64(monthEnd.Day()), now,
			)
			if err != nil {
				return fmt.Errorf("failed to insert resource usage history: %w", err)
			}
		}
	}

	// Create billing periods for the last 3 months
	for i := 0; i < 3; i++ {
		monthStart := time.Date(now.Year(), now.Month()-time.Month(i), 1, 0, 0, 0, 0, time.UTC)
		monthEnd := monthStart.AddDate(0, 1, -1)

		// Calculate monthly costs by resource type
		computeCost := 0.0
		storageCost := 0.0
		lambdaCost := 0.0
		rdbCost := 0.0
		secretCost := 0.0

		// Query historical data to calculate costs
		rows, err := db.Query(`
			SELECT resource_type, SUM(total_cost) as total_cost
			FROM resource_usage_histories 
			WHERE organization_id = ? AND period_start >= ? AND period_end <= ?
			GROUP BY resource_type
		`, orgID, monthStart, monthEnd)
		if err != nil {
			return fmt.Errorf("failed to query historical data: %w", err)
		}

		for rows.Next() {
			var resourceType string
			var cost float64
			if err := rows.Scan(&resourceType, &cost); err != nil {
				rows.Close()
				return fmt.Errorf("failed to scan historical data: %w", err)
			}

			switch resourceType {
			case string(models.ResourceTypeCompute):
				computeCost = cost
			case string(models.ResourceTypeStorage):
				storageCost = cost
			case string(models.ResourceTypeLambda):
				lambdaCost = cost
			case string(models.ResourceTypeRDB):
				rdbCost = cost
			case string(models.ResourceTypeSecret):
				secretCost = cost
			}
		}
		rows.Close()

		totalCost := computeCost + storageCost + lambdaCost + rdbCost + secretCost
		isPaid := i > 0 // Only the current month is unpaid

		_, err = db.Exec(`
			INSERT OR REPLACE INTO billing_periods (
				id, organization_id, period_start, period_end, total_cost, currency, is_paid,
				compute_cost, lambda_cost, storage_cost, rdb_cost, secret_cost, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			uuid.New().String(), orgID, monthStart, monthEnd, totalCost, "USD", isPaid,
			computeCost, lambdaCost, storageCost, rdbCost, secretCost, now, now,
		)
		if err != nil {
			return fmt.Errorf("failed to insert billing period: %w", err)
		}
	}

	return nil
}
