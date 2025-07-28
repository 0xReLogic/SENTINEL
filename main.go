package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/sentinel/checker"
	"github.com/sentinel/config"
)

func main() {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}

	// Construct the path to the config file
	configPath := filepath.Join(cwd, "sentinel.yaml")

	// Load the configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	fmt.Println("SENTINEL Monitoring System")
	fmt.Printf("Loaded %d services to monitor\n", len(cfg.Services))
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println("-----------------------------------")

	// Create a ticker that triggers every minute
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// Run the first check immediately
	runChecks(cfg)

	// Then run on ticker schedule
	for range ticker.C {
		runChecks(cfg)
	}
}

// runChecks performs checks on all services in the configuration
func runChecks(cfg *config.Config) {
	fmt.Printf("\n[%s] Running service checks...\n", time.Now().Format("2006-01-02 15:04:05"))

	for _, service := range cfg.Services {
		// Check the service
		status := checker.CheckService(service.Name, service.URL)

		// Print the result
		fmt.Println(status)
	}

	fmt.Println("-----------------------------------")
}
