// Package cmd implements the command-line interface for SENTINEL
package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/0xReLogic/SENTINEL/config"
	"github.com/spf13/cobra"
)

var (
	// configPath holds the path to the configuration file
	configPath string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   appName,
	Short: descShort,
	Long:  fmt.Sprintf(descLong, appRepository),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitError)
	}
}

func init() {
	// add persistent flags that are available to all subcommands
	rootCmd.PersistentFlags().StringVarP(&configPath, flagConfig, flagConfigShort,
		defaultConfigFile, descConfigFlag)
}

// loadConfig loads configuration from the specified path with helpful error messages
func loadConfig() (*config.Config, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf(errInvalidConfigPath, err)
	}

	cfg, err := config.LoadConfig(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf(errConfigNotFound, configPath, defaultConfigFile, flagConfig)
		}
		return nil, err
	}

	return cfg, nil
}

// validateServices validates all services in the configuration
func validateServices(services []config.Service) []error {
	var errors []error

	for i, service := range services {
		if service.Name == "" {
			errors = append(errors,
				fmt.Errorf(errServiceNameReq, i+1))
		}
		if service.URL == "" {
			errors = append(errors,
				fmt.Errorf(errServiceURLReq, i+1, service.Name))
		}
		// validate URL format if provided
		if service.URL != "" && !isValidURL(service.URL) {
			errors = append(errors,
				fmt.Errorf(errServiceURLInvalid, i+1, service.Name, service.URL))
		}
	}

	return errors
}

// isValidURL checks if a string is a valid HTTP/HTTPS URL
func isValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Host != "" &&
		(u.Scheme == schemeHTTP || u.Scheme == schemeHTTPS)
}
