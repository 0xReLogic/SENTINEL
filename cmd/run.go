package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/0xReLogic/SENTINEL/checker"
	"github.com/0xReLogic/SENTINEL/config"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   cmdNameRun,
	Short: descRunShort,
	Long:  descRunLong,
	Run: func(cmd *cobra.Command, args []string) {
		// load configuration
		cfg, err := loadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, errLoadingConfig, err)
			os.Exit(exitConfigError)
		}

		// print startup banner
		printBanner(cfg)

		// create a ticker that triggers at the configured interval
		ticker := time.NewTicker(checkInterval)
		defer ticker.Stop()

		// run the first check immediately
		runChecks(cfg)

		// then run on ticker schedule
		for range ticker.C {
			runChecks(cfg)
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

// runChecks performs checks on all services in the configuration
func runChecks(cfg *config.Config) {
	fmt.Printf(fmtTimestamp, time.Now().Format(timestampFormat), msgRunningChecks)

	for _, service := range cfg.Services {
		status := checker.CheckService(service.Name, service.URL)
		fmt.Println(status)
	}

	fmt.Println(separator)
}
