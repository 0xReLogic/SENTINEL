package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sentinel/checker"
	"github.com/sentinel/config"
	"github.com/spf13/cobra"
)

var (
	configFile string
	cfg        *config.Config
)

func main() {
	// Define the root command
	rootCmd := &cobra.Command{
		Use:   "sentinel",
		Short: "SENTINEL - Simple service monitoring tool",
		Long:  `SENTINEL is a simple monitoring tool that checks the status of web services and reports their status.`,
	}

	// Add global flags
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "sentinel.yaml", "Path to configuration file")

	// Add run command
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Start monitoring services",
		Long:  `Start monitoring services defined in the configuration file.`,
		RunE:  runCommand,
	}

	// Add validate command
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate configuration file",
		Long:  `Validate the configuration file without starting the monitoring.`,
		RunE:  validateCommand,
	}

	// Add once command
	onceCmd := &cobra.Command{
		Use:   "once",
		Short: "Run checks once and exit",
		Long:  `Run checks on all services once and exit without continuous monitoring.`,
		RunE:  onceCommand,
	}

	// Add commands to root
	rootCmd.AddCommand(runCmd, validateCmd, onceCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// loadConfig loads the configuration file
func loadConfig() error {
	// Get absolute path if not already
	if !filepath.IsAbs(configFile) {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting current directory: %w", err)
		}
		configFile = filepath.Join(cwd, configFile)
	}

	// Load the configuration
	var err error
	cfg, err = config.LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	return nil
}

// validateCommand validates the configuration file
func validateCommand(cmd *cobra.Command, args []string) error {
	if err := loadConfig(); err != nil {
		return err
	}

	fmt.Println("Configuration file is valid!")
	fmt.Printf("Found %d services to monitor\n", len(cfg.Services))
	fmt.Printf("Check interval: %s\n", cfg.Settings.CheckInterval)
	fmt.Printf("Default timeout: %s\n", cfg.Settings.DefaultTimeout)

	return nil
}

// onceCommand runs checks once and exits
func onceCommand(cmd *cobra.Command, args []string) error {
	if err := loadConfig(); err != nil {
		return err
	}

	fmt.Println("SENTINEL Monitoring System - One-time Check")
	fmt.Printf("Checking %d services...\n", len(cfg.Services))
	fmt.Println("-----------------------------------")

	runChecks(cfg)
	return nil
}

// runCommand starts the continuous monitoring
func runCommand(cmd *cobra.Command, args []string) error {
	if err := loadConfig(); err != nil {
		return err
	}

	fmt.Println("SENTINEL Monitoring System")
	fmt.Printf("Loaded %d services to monitor\n", len(cfg.Services))
	fmt.Printf("Check interval: %s\n", cfg.Settings.CheckInterval)
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println("-----------------------------------")

	// Create a ticker that triggers at the configured interval
	ticker := time.NewTicker(cfg.Settings.CheckInterval)
	defer ticker.Stop()

	// Run the first check immediately
	runChecks(cfg)

	// Then run on ticker schedule
	for range ticker.C {
		runChecks(cfg)
	}

	return nil
}

// runChecks performs checks on all services in the configuration concurrently
func runChecks(cfg *config.Config) {
	fmt.Printf("\n[%s] Running service checks...\n", time.Now().Format("2006-01-02 15:04:05"))

	// Create a channel to collect results
	results := make(chan checker.ServiceStatus, len(cfg.Services))

	// Start a goroutine for each service
	for _, service := range cfg.Services {
		go func(svc config.Service) {
			// Get the appropriate timeout for this service
			timeout := cfg.GetServiceTimeout(svc)

			// Check the service and send result to channel
			results <- checker.CheckService(svc.Name, svc.URL, timeout)
		}(service)
	}

	// Collect and print results
	for i := 0; i < len(cfg.Services); i++ {
		status := <-results
		fmt.Println(status)
	}

	fmt.Println("-----------------------------------")
}
