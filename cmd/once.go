package cmd

import (
	"fmt"
	"os"
	"time"
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

    printBanner(cfg)

    // Create StateManager once here
    stateManager := NewStateManager()
    
    ticker := time.NewTicker(checkInterval)
    defer ticker.Stop()

    // Pass stateManager to function
    runChecksAndGetStatus(cfg, stateManager)

    for range ticker.C {
        runChecksAndGetStatus(cfg, stateManager)
    }
},
}

func init() {
	rootCmd.AddCommand(onceCmd)
}
