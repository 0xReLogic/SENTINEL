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
	"github.com/0xReLogic/SENTINEL/storage"
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
// NotificationConfig interface for generic notification handling
type NotificationConfig interface {
	IsEnabled() bool
	GetNotifyOn() []string
}

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

func (sm *StateManager) ProcessStatus(status checker.ServiceStatus, service config.Service, cfg config.TelegramConfig) NotificationAction {
	checkTime := time.Now()
	previousIsUp, exists := sm.serviceState[status.URL]
	// Defer the state update so it happens regardless of how the function exits.
	defer func() { sm.serviceState[status.URL] = status.IsUp }()

	// 1. Determine the exact state transition
	isDownTransition := exists && previousIsUp && !status.IsUp
	isRecoveryTransition := exists && !previousIsUp && status.IsUp

	// 2. If there's no state change, we are done.
	if !isDownTransition && !isRecoveryTransition {
		return NotificationAction{Action: NoAction}
	}
	// 3. Check for throttling. If a notification was sent recently, we are done.
	isThrottled := time.Since(sm.lastNotificationTime[status.URL]) < service.Interval
	if isThrottled {
		return NotificationAction{Action: NoAction}
	}
	// Handle a RECOVERY notification
	if isRecoveryTransition && contains(cfg.NotifyOn, "recovery") {
		sm.lastNotificationTime[status.URL] = checkTime
		downtime := time.Since(sm.serviceDownSince[status.URL])
		delete(sm.serviceDownSince, status.URL)
		return NotificationAction{Action: NotifyRecovery, Downtime: downtime}
	}
	// Handle a DOWN notification
	if isDownTransition && contains(cfg.NotifyOn, "down") {
		sm.lastNotificationTime[status.URL] = checkTime
		sm.serviceDownSince[status.URL] = checkTime
		return NotificationAction{Action: NotifyDown}
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

// NotifyDiscordServiceDown sends a Discord notification when a service goes DOWN
func NotifyDiscordServiceDown(cfg config.DiscordConfig, status checker.ServiceStatus, checkTime time.Time) {
	var errorMsg string
	if status.Error != nil {
		errorMsg = status.Error.Error()
	} else {
		errorMsg = fmt.Sprintf("HTTP Status Code %d", status.StatusCode)
	}
	embed := notifier.FormatDownEmbed(status.Name, status.URL, errorMsg, checkTime)

	log.Printf("INFO: Sending Discord DOWN notification for %s", status.Name)

	err := notifier.SendDiscordNotification(cfg.WebhookURL, "", embed)
	if err != nil {
		log.Printf("ERROR: Failed to send Discord DOWN notification for %s: %v", status.Name, err)
	}
}

// NotifyDiscordServiceRecovery sends a Discord notification when a service RECOVERS
func NotifyDiscordServiceRecovery(cfg config.DiscordConfig, status checker.ServiceStatus, downtime time.Duration, recoveryTime time.Time) {
	embed := notifier.FormatRecoveryEmbed(status.Name, status.URL, downtime, recoveryTime)

	log.Printf("INFO: Sending Discord RECOVERY notification for %s", status.Name)

	err := notifier.SendDiscordNotification(cfg.WebhookURL, "", embed)
	if err != nil {
		log.Printf("ERROR: Failed to send Discord RECOVERY notification for %s: %v", status.Name, err)
	}
}

// processNotifications handles both Telegram and Discord notifications for a service status
func processNotifications(cfg *config.Config, stateManager *StateManager, status checker.ServiceStatus, service config.Service) {
	// Process Telegram notifications
	if cfg.Notifications.Telegram.Enabled {
		action := stateManager.ProcessStatus(status, service, cfg.Notifications.Telegram)
		switch action.Action {
		case NotifyDown:
			log.Printf("INFO: Service '%s' is DOWN. Preparing Telegram notification.", status.Name)
			NotifyServiceDown(cfg.Notifications.Telegram, status, time.Now())
		case NotifyRecovery:
			log.Printf("INFO: Service '%s' has RECOVERED. Preparing Telegram notification.", status.Name)
			NotifyServiceRecovery(cfg.Notifications.Telegram, status, action.Downtime, time.Now())
		}
	}

	// Process Discord notifications
	if cfg.Notifications.Discord.Enabled {
		tempCfg := config.TelegramConfig{
			Enabled:  cfg.Notifications.Discord.Enabled,
			NotifyOn: cfg.Notifications.Discord.NotifyOn,
		}
		action := stateManager.ProcessStatus(status, service, tempCfg)
		switch action.Action {
		case NotifyDown:
			log.Printf("INFO: Service '%s' is DOWN. Preparing Discord notification.", status.Name)
			NotifyDiscordServiceDown(cfg.Notifications.Discord, status, time.Now())
		case NotifyRecovery:
			log.Printf("INFO: Service '%s' has RECOVERED. Preparing Discord notification.", status.Name)
			NotifyDiscordServiceRecovery(cfg.Notifications.Discord, status, action.Downtime, time.Now())
		}
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

func runChecksAndGetStatus(cfg *config.Config, stateManager *StateManager, store storage.Storage) bool {
	fmt.Printf("[%s] --- Running Checks ---\n", time.Now().Format("2006-01-02 15:04:05"))
	allUp := true

	for _, service := range cfg.Services {
		status := checker.CheckService(service.Name, service.URL, service.Timeout)
		fmt.Println(status)

		if !status.IsUp {
			allUp = false
		}

		// Save to storage if configured
		if store != nil {
			if err := store.SaveCheck(status); err != nil {
				log.Printf("ERROR: Failed to save check to storage: %v", err)
			}
		}

		processNotifications(cfg, stateManager, status, service)
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