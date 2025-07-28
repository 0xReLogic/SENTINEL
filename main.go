package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sentinel/checker"
	"github.com/sentinel/config"
	"github.com/sentinel/notifier"
	"github.com/sentinel/storage"
	"github.com/spf13/cobra"
)

var (
	configFile string
	cfg        *config.Config
	store      *storage.StatusStorage
	notifiers  []notifier.Notifier
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

// loadConfig loads the configuration file and initializes required components
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

	// Initialize status storage
	store, err = storage.NewStatusStorage(cfg.Settings.StoragePath)
	if err != nil {
		return fmt.Errorf("error initializing status storage: %w", err)
	}

	// Initialize notifiers
	notifiers = cfg.GetNotifiers()
	if len(notifiers) > 0 {
		fmt.Printf("Initialized %d notification channels\n", len(notifiers))
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
			var result checker.ServiceStatus

			// Handle different check types
			if svc.CheckType == "tcp" && svc.Port > 0 {
				// Use TCP port check
				result = checker.CheckTCPPort(svc.Name, svc.URL, svc.Port, timeout)
			} else if svc.ExpectString != "" {
				// Use content check if expect_string is configured
				result = checker.CheckContent(svc.Name, svc.URL, svc.ExpectString, timeout)
			} else {
				// Default to HTTP check
				result = checker.CheckService(svc.Name, svc.URL, timeout)

				// Check for expected status code if configured
				if svc.ExpectStatus > 0 && result.IsUp {
					if result.StatusCode != svc.ExpectStatus {
						result.IsUp = false
						result.Error = fmt.Errorf("expected status code %d, got %d", svc.ExpectStatus, result.StatusCode)
					}
				}
			}

			results <- result
		}(service)
	}

	// Collect results and check for status changes
	for i := 0; i < len(cfg.Services); i++ {
		status := <-results
		fmt.Println(status)

		// Check if status has changed and send notifications if needed
		if store != nil && len(notifiers) > 0 {
			changed, oldStatus := store.HasStatusChanged(status.Name, status)

			// If this is a new service or status has changed
			if changed {
				var notificationType notifier.NotificationType

				if status.IsUp {
					notificationType = notifier.ServiceRecovered
				} else {
					notificationType = notifier.ServiceDown
				}

				// Create notification
				notification := notifier.Notification{
					Type:        notificationType,
					ServiceName: status.Name,
					NewStatus:   status,
					Timestamp:   time.Now(),
				}

				// Add old status if it exists
				if oldStatus.Name != "" {
					notification.OldStatus = &oldStatus
				}

				// Send notification to all configured notifiers
				for _, n := range notifiers {
					go func(notif notifier.Notifier, note notifier.Notification) {
						if err := notif.SendNotification(note); err != nil {
							fmt.Printf("Error sending notification: %v\n", err)
						}
					}(n, notification)
				}
			}
		}

		// Store the current status
		if store != nil {
			store.SetStatus(status)
		}
	}

	fmt.Println("-----------------------------------")
}
