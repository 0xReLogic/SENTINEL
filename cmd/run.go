package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/0xReLogic/SENTINEL/config"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   cmdNameRun,
	Short: descRunShort,
	Long:  descRunLong,
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
	rootCmd.AddCommand(runCmd)
}

// printBanner displays the startup banner
func printBanner(cfg *config.Config) {
	fmt.Println(bannerTitle)
	fmt.Printf(fmtLoadedServices, len(cfg.Services))
	fmt.Println(bannerExitInstruction)
	fmt.Println(separator)
}
