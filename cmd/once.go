package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/0xReLogic/SENTINEL/config"
	"github.com/spf13/cobra"
)

// onceCmd represents the once command
var onceCmd = &cobra.Command{
	Use:   cmdNameOnce,
	Short: descOnceShort,
	Long:  fmt.Sprintf(descOnceLong, exitSuccess, exitError, exitConfigError),
	Run: func(cmd *cobra.Command, args []string) {
		// load configuration
		cfg, err := loadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, errLoadingConfig, err)
			os.Exit(exitConfigError)
		}

		// run checks once
		fmt.Printf(fmtTimestamp, time.Now().Format(timestampFormat), msgRunningChecks)

		allUp := runChecksWithStatus(cfg)

		fmt.Println(separator)

		// exit with appropriate code
		if !allUp {
			os.Exit(exitError)
		}
		os.Exit(exitSuccess)
	},
}

func init() {
	rootCmd.AddCommand(onceCmd)
}

// runChecksWithStatus performs checks and returns overall status
func runChecksWithStatus(cfg *config.Config) bool {
	return runChecksAndGetStatus(cfg)
}
