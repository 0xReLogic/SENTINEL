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

type ActionType int

const (
	// NoAction means no notification should be sent.
	// NotifyDown means a "service down" notification should be sent.
	// NotifyRecovery means a "service recovered" notification should be sent.
	NoAction ActionType = iota
	NotifyDown
	NotifyRecovery
)

// NotificationAction represents the decision made by the StateManager about whether a notification should be sent.
// It includes the type of action and any additional data needed, like the total downtime for a recovery alert.
type NotificationAction struct {
	Action   ActionType
	Downtime time.Duration
}

// StateManager encapsulates the state and logic for tracking service statuses over time.
type StateManager struct {
	serviceState         map[string]bool
	lastNotificationTime map[string]time.Time
	serviceDownSince     map[string]time.Time
}

// NewStateManager creates and initializes a new StateManager.
func NewStateManager() *StateManager {
	return &StateManager{
		serviceState:         make(map[string]bool),
		lastNotificationTime: make(map[string]time.Time),
		serviceDownSince:     make(map[string]time.Time),
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   appName,
	Short: descShort,
	Long:  fmt.Sprintf(descLong, appRepository),
}

// ProcessStatus checks the current status of a service against its historical state and determines if a notification action is required ONLY on state transitions.
// CHANGED: Added 'service config.Service' parameter to get the correct interval.
func (sm *StateManager) ProcessStatus(status checker.ServiceStatus, service config.Service, cfg config.TelegramConfig) NotificationAction {
	checkTime := time.Now()
	previousIsUp, exists := sm.serviceState[status.URL]
	isDownTransition := exists && previousIsUp && !status.IsUp
	isRecoveryTransition := exists && !previousIsUp && status.IsUp
	sm.serviceState[status.URL] = status.IsUp
	if isDownTransition {
		sm.serviceDownSince[status.URL] = checkTime
	}

	// Determine if a notification should be sent based ONLY on transitions.
	if isRecoveryTransition {
		// CHANGED: Replaced 'checkInterval' with 'service.Interval' for accurate throttling.
		if contains(cfg.NotifyOn, "recovery") && time.Since(sm.lastNotificationTime[status.URL]) >= service.Interval {
			sm.lastNotificationTime[status.URL] = checkTime
			downtime := time.Since(sm.serviceDownSince[status.URL])
			delete(sm.serviceDownSince, status.URL)
			return NotificationAction{Action: NotifyRecovery, Downtime: downtime}
		}

	} else if isDownTransition {
		// CHANGED: Replaced 'checkInterval' with 'service.Interval' for accurate throttling.
		if contains(cfg.NotifyOn, "down") && time.Since(sm.lastNotificationTime[status.URL]) >= service.Interval {
			sm.lastNotificationTime[status.URL] = checkTime
			return NotificationAction{Action: NotifyDown}
		}
	}

	return NotificationAction{Action: NoAction}
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
func NotifyServiceDown(cfg config.TelegramConfig, status checker.ServiceStatus, checkTime time.Time) {
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
		log.Printf("ERROR: Failed to send Telegram DOWN notification for %s: %v", status.Name, err)
	}
}

// notifyServiceRecovery constructs and sends a 'Service RECOVERED' notification to Telegram.
func NotifyServiceRecovery(cfg config.TelegramConfig, status checker.ServiceStatus, downtime time.Duration, recoveryTime time.Time) {
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
		if service.Interval <= 0 {
			errors = append(errors,
				fmt.Errorf(errServiceIntervalInvalid, i+1, service.Name))
		}
		if service.Timeout <= 0 {
			errors = append(errors,
				fmt.Errorf(errServiceTimeoutInvalid, i+1, service.Name))
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

func runChecksAndGetStatus(cfg *config.Config, stateManager *StateManager) bool {
	fmt.Printf("[%s] --- Running Checks ---\n", time.Now().Format("2006-01-02 15:04:05"))
	allUp := true
	telegramEnabled := cfg.Notifications.Telegram.Enabled

	for _, service := range cfg.Services {

		status := checker.CheckService(service.Name, service.URL, service.Timeout)
		fmt.Println(status)

		if !status.IsUp {
			allUp = false
		}

		if telegramEnabled {
			// CHANGED: Passed the 'service' object to the function call.
			action := stateManager.ProcessStatus(status, service, cfg.Notifications.Telegram)

			switch action.Action {
			case NotifyDown:
				log.Printf("INFO: Service '%s' is DOWN. Preparing notification.", status.Name)
				NotifyServiceDown(cfg.Notifications.Telegram, status, time.Now())
			case NotifyRecovery:
				log.Printf("INFO: Service '%s' has RECOVERED. Preparing notification.", status.Name)
				NotifyServiceRecovery(cfg.Notifications.Telegram, status, action.Downtime, time.Now())
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