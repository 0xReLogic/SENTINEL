// Package cmd implements the command-line interface for SENTINEL
package cmd

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/0xReLogic/SENTINEL/checker"
	"github.com/0xReLogic/SENTINEL/config"
	"github.com/0xReLogic/SENTINEL/notifier"
	"github.com/spf13/cobra"
)

var (
	// configPath holds the path to the configuration file, set by a command-line flag.
	configPath string
)

var (
	// configPath holds the path to the configuration file, set by a command-line flag.
	// lastNotificationTime records the timestamp of the last alert sent for each service to enable rate limiting.
	// serviceDownSince records the timestamp when a service first went down to calculate total downtime upon recovery.
	serviceState         = make(map[string]bool)
	lastNotificationTime = make(map[string]time.Time)
	serviceDownSince     = make(map[string]time.Time)
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
func loadConfig(path string) (*config.Config, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf(errInvalidConfigPath, err)
	}

	cfg, err := config.LoadConfig(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf(errConfigNotFound, path, defaultConfigFile, flagConfig)
		}
		return nil, err
	}

	return cfg, nil
}

// notifyServiceDown constructs and sends a 'Service DOWN' notification to Telegram.
func notifyServiceDown(cfg config.TelegramConfig, status checker.ServiceStatus, checkTime time.Time) {
	var errorMsg string
	if status.Error != nil {
		errorMsg = status.Error.Error()
	} else {
		errorMsg = fmt.Sprintf("HTTP Status Code %d", status.StatusCode)
	}
	message := notifier.FormatDownMessage(status.Name, status.URL, errorMsg, checkTime)

	log.Printf("INFO: Sending DOWN notification for %s", status.Name)

	err := notifier.SendTelegramNotification(cfg.BotToken, cfg.ChatID, message)
	if err != nil {
		log.Printf("ERROR: Failed to send Telegram DOWN notification for %s", status.Name)
	}
}

// notifyServiceRecovery constructs and sends a 'Service RECOVERED' notification to Telegram.
func notifyServiceRecovery(cfg config.TelegramConfig, status checker.ServiceStatus, downtime time.Duration, recoveryTime time.Time) {
	message := notifier.FormatRecoveryMessage(status.Name, status.URL, downtime, recoveryTime)

	log.Printf("INFO: Sending RECOVERY notification for %s", status.Name)

	err := notifier.SendTelegramNotification(cfg.BotToken, cfg.ChatID, message)
	if err != nil {
		log.Printf("ERROR: Failed to send Telegram RECOVERY notification for %s: %v", status.Name, err)
	}
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

// runChecksAndGetStatus iterates through services, checks their status, and triggers notifications on state changes.
func runChecksAndGetStatus(cfg *config.Config) bool {

	fmt.Printf("[%s] --- Running Checks ---\n", time.Now().Format("2006-01-02 15:04:05"))
	allUp := true
	checkTime := time.Now()
	telegramEnabled := cfg.Notifications.Telegram.Enabled

	for _, service := range cfg.Services {

		status := checker.CheckService(service.Name, service.URL)
		fmt.Println(status)
		if !status.IsUp {
			allUp = false
		}

		previousIsUp, exists := serviceState[service.URL]

		isDownTransition := (!exists && !status.IsUp) || (exists && previousIsUp && !status.IsUp)
		isRecoveryTransition := exists && !previousIsUp && status.IsUp

		serviceState[service.URL] = status.IsUp
		if isDownTransition {
			serviceDownSince[service.URL] = checkTime
		}
		if telegramEnabled {

			if isRecoveryTransition {
				if contains(cfg.Notifications.Telegram.NotifyOn, "recovery") && time.Since(lastNotificationTime[service.URL]) >= checkInterval {
					downtime := time.Since(serviceDownSince[service.URL])
					notifyServiceRecovery(cfg.Notifications.Telegram, status, downtime, checkTime)
					lastNotificationTime[service.URL] = checkTime
					delete(serviceDownSince, service.URL)
				}
			} else if !status.IsUp {
				if contains(cfg.Notifications.Telegram.NotifyOn, "down") && time.Since(lastNotificationTime[service.URL]) >= checkInterval {
					log.Printf("INFO: Service '%s' is DOWN. Preparing notification.", status.Name)
					notifyServiceDown(cfg.Notifications.Telegram, status, checkTime)
					lastNotificationTime[service.URL] = checkTime
				}
			}
		}
	}

	fmt.Println("---------------------------------------")
	return allUp
}

// contains checks if a string exists within a slice of strings.
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
