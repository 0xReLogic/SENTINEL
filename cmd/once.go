// in cmd/once.go

package cmd

import (
	"fmt"
	"os"

	"github.com/0xReLogic/SENTINEL/storage"
	"github.com/spf13/cobra"
)

// onceCmd represents the once command
var onceCmd = &cobra.Command{
	Use:   cmdNameOnce,
	Short: descOnceShort,
	Long:  fmt.Sprintf(descOnceLong, exitSuccess, exitError, exitConfigError),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := loadConfig(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, errLoadingConfig, err)
			os.Exit(exitConfigError)
		}

		// Initialize storage if configured
		var store storage.Storage
		if cfg.Storage.Type == "sqlite" && cfg.Storage.Path != "" {
			store, err = storage.NewSQLiteStorage(cfg.Storage.Path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to initialize storage: %v\n", err)
			} else {
				defer store.Close()
			}
		}

		// Create StateManager to handle notifications correctly even on a single run.
		stateManager := NewStateManager()
		allServicesUp := runChecksAndGetStatus(cfg, stateManager, store)

		// Exit with the correct status code based on the result.
		if allServicesUp {
			fmt.Println("\nAll services are UP.")
			os.Exit(exitSuccess) // Exit with 0
		} else {
			fmt.Println("\nOne or more services are DOWN.")
			os.Exit(exitError)   // Exit with 1
		}
	},
}

func init() {
	rootCmd.AddCommand(onceCmd)
}