package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   cmdNameValidate,
	Short: descValidShort,
	Long:  descValidLong,
	Run: func(cmd *cobra.Command, args []string) {
		// load configuration
		cfg, err := loadConfig()
		if err != nil {
			fmt.Fprint(os.Stderr, msgValidationFailed)
			fmt.Fprintf(os.Stderr, indent+"%v\n", err)
			os.Exit(exitConfigError)
		}

		// validate that services exist
		if len(cfg.Services) == 0 {
			fmt.Fprint(os.Stderr, msgValidationFailed)
			fmt.Fprintln(os.Stderr, indent+msgNoServicesDefined)
			os.Exit(exitConfigError)
		}

		// validate each service
		errors := validateServices(cfg.Services)
		if len(errors) > 0 {
			fmt.Fprint(os.Stderr, msgValidationFailed)
			for _, err := range errors {
				fmt.Fprintf(os.Stderr, listPrefix+"%v\n", err)
			}
			os.Exit(exitConfigError)
		}

		// validation successful
		fmt.Println(msgValidationSuccess)
		fmt.Printf(fmtLoadedServicesValidation, len(cfg.Services))
		fmt.Println(msgServicesConfigured)
		for i, service := range cfg.Services {
			fmt.Printf(fmtServiceListItem, i+1, service.Name, service.URL)
		}
		os.Exit(exitSuccess)
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
