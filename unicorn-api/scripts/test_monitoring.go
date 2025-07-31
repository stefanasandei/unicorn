package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"unicorn-api/internal/stores"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Open SQLite database
	db, err := sql.Open("sqlite3", "./unicorn.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Create monitoring store
	monitoringStore, err := stores.NewGORMMonitoringStore("./unicorn.db")
	if err != nil {
		log.Fatal("Failed to create monitoring store:", err)
	}

	// Test organization ID (use one from the database)
	orgID := "4bca97d4-81ac-4a20-9db4-3799afdca20c"

	// Test getting active resources
	fmt.Println("Testing GetActiveResourcesByOrganization...")
	activeResources, err := monitoringStore.GetActiveResourcesByOrganization(orgID)
	if err != nil {
		log.Fatal("Failed to get active resources:", err)
	}

	fmt.Printf("Found %d active resources\n", len(activeResources))
	for i, resource := range activeResources {
		fmt.Printf("Resource %d:\n", i+1)
		fmt.Printf("  ID: %s\n", resource.ResourceID)
		fmt.Printf("  Name: %s\n", resource.ResourceName)
		fmt.Printf("  Type: %s\n", resource.ResourceType)
		fmt.Printf("  Status: %s\n", resource.Status)
		fmt.Printf("  CPU Usage: %.2f%%\n", resource.CPUUsage)
		fmt.Printf("  Memory Usage: %d MB\n", resource.MemoryUsage)
		fmt.Printf("  Storage Usage: %d MB\n", resource.StorageUsage)
		fmt.Printf("  Network Usage: %d MB\n", resource.NetworkUsage)
		fmt.Printf("  Cost Per Hour: $%.2f\n", resource.CostPerHour)
		fmt.Printf("  Total Cost: $%.2f\n", resource.TotalCost)
		fmt.Printf("  Currency: %s\n", resource.Currency)
		fmt.Printf("  Created: %s\n", resource.ResourceCreatedAt.Format("2006-01-02 15:04:05"))
		if resource.LastActiveAt != nil {
			fmt.Printf("  Last Active: %s\n", resource.LastActiveAt.Format("2006-01-02 15:04:05"))
		}
		fmt.Println()
	}

	// Test getting resource usage summary
	fmt.Println("Testing GetResourceUsageSummary...")
	start := time.Now().AddDate(0, 0, -30) // 30 days ago
	end := time.Now()
	summary, err := monitoringStore.GetResourceUsageSummary(orgID, start, end)
	if err != nil {
		log.Fatal("Failed to get resource usage summary:", err)
	}

	fmt.Printf("Resource Usage Summary:\n")
	fmt.Printf("  Total Resources: %d\n", summary.TotalResources)
	fmt.Printf("  Active Resources: %d\n", summary.ActiveResources)
	fmt.Printf("  Total Cost: $%.2f\n", summary.TotalCost)
	fmt.Printf("  Currency: %s\n", summary.Currency)
	fmt.Println()

	// Test getting billing history
	fmt.Println("Testing GetBillingPeriodsByOrganization...")
	billingPeriods, err := monitoringStore.GetBillingPeriodsByOrganization(orgID)
	if err != nil {
		log.Fatal("Failed to get billing periods:", err)
	}

	fmt.Printf("Found %d billing periods\n", len(billingPeriods))
	for i, period := range billingPeriods {
		fmt.Printf("Billing Period %d:\n", i+1)
		fmt.Printf("  Period: %s to %s\n",
			period.PeriodStart.Format("2006-01-02"),
			period.PeriodEnd.Format("2006-01-02"))
		fmt.Printf("  Total Cost: $%.2f\n", period.TotalCost)
		fmt.Printf("  Compute Cost: $%.2f\n", period.ComputeCost)
		fmt.Printf("  Storage Cost: $%.2f\n", period.StorageCost)
		fmt.Printf("  Lambda Cost: $%.2f\n", period.LambdaCost)
		fmt.Printf("  RDB Cost: $%.2f\n", period.RDBCost)
		fmt.Printf("  Secret Cost: $%.2f\n", period.SecretCost)
		fmt.Printf("  Is Paid: %t\n", period.IsPaid)
		fmt.Println()
	}

	fmt.Println("✅ Monitoring service test completed successfully!")
}
