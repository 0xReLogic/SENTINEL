package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// onceCmd represents the once command
var onceCmd = &cobra.Command{
	Use:   cmdNameOnce,
	Short: descOnceShort,
	Long:  fmt.Sprintf(descOnceLong, exitSuccess, exitError, exitConfigError),
	Run: func(cmd *cobra.Command, args []string) {
		// load configuration
		cfg, err := loadConfig(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, errLoadingConfig, err)
			os.Exit(exitConfigError)
		}

		// run checks once
		allUp := runChecksAndGetStatus(cfg)

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
