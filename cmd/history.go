package cmd

import (
	"fmt"
	"os"

	"github.com/0xReLogic/SENTINEL/storage"
	"github.com/spf13/cobra"
)

var historyLimit int

var historyCmd = &cobra.Command{
	Use:   "history [service-name]",
	Short: "View check history for a service",
	Long:  "Display historical check results for a specific service from the database",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceName := args[0]

		cfg, err := loadConfig(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(exitConfigError)
		}

		if cfg.Storage.Type == "" || cfg.Storage.Path == "" {
			fmt.Fprintln(os.Stderr, "Storage not configured. Please configure storage in sentinel.yaml")
			os.Exit(exitConfigError)
		}

		store, err := storage.NewSQLiteStorage(cfg.Storage.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening storage: %v\n", err)
			os.Exit(exitError)
		}
		defer store.Close()

		records, err := store.GetHistory(serviceName, historyLimit)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error retrieving history: %v\n", err)
			os.Exit(exitError)
		}

		if len(records) == 0 {
			fmt.Printf("No history found for service: %s\n", serviceName)
			return
		}

		fmt.Printf("Check History for '%s' (last %d records):\n\n", serviceName, len(records))
		fmt.Println("Time                 | Status | Response Time | Status Code | Error")
		fmt.Println("---------------------|--------|---------------|-------------|-------")

		for _, record := range records {
			status := "UP  "
			if !record.IsUp {
				status = "DOWN"
			}

			responseTime := fmt.Sprintf("%dms", record.ResponseTimeMs)
			statusCode := fmt.Sprintf("%d", record.StatusCode)
			if record.StatusCode == 0 {
				statusCode = "N/A"
			}

			errorMsg := record.ErrorMessage
			if errorMsg == "" {
				errorMsg = "-"
			} else if len(errorMsg) > 40 {
				errorMsg = errorMsg[:37] + "..."
			}

			fmt.Printf("%s | %s   | %-13s | %-11s | %s\n",
				record.CheckedAt.Format("2006-01-02 15:04:05"),
				status,
				responseTime,
				statusCode,
				errorMsg,
			)
		}
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
	historyCmd.Flags().IntVarP(&historyLimit, "limit", "l", 50, "Number of records to display")
}
