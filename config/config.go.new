package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sentinel/notifier"
	"gopkg.in/yaml.v3"
)

// Service represents a single service to be monitored
type Service struct {
	Name          string        `yaml:"name"`
	URL           string        `yaml:"url"`
	Timeout       time.Duration `yaml:"timeout,omitempty"`       // Optional timeout per service
	ExpectString  string        `yaml:"expect_string,omitempty"` // String to expect in the response
	ExpectStatus  int           `yaml:"expect_status,omitempty"` // Expected HTTP status code
	CheckType     string        `yaml:"check_type,omitempty"`    // Type of check: http, tcp, etc.
	Port          int           `yaml:"port,omitempty"`          // Port for TCP checks
}

// NotificationConfig represents notification settings
type NotificationConfig struct {
	Discord *notifier.WebhookConfig `yaml:"discord,omitempty"` // Discord webhook configuration
	Slack   *notifier.WebhookConfig `yaml:"slack,omitempty"`   // Slack webhook configuration
	Webhook *notifier.WebhookConfig `yaml:"webhook,omitempty"` // Generic webhook configuration
}

// Settings represents global application settings
type Settings struct {
	CheckInterval  time.Duration `yaml:"check_interval,omitempty"`  // Interval between checks
	DefaultTimeout time.Duration `yaml:"default_timeout,omitempty"` // Default timeout for all services
	StoragePath    string        `yaml:"storage_path,omitempty"`    // Path to store status data
}

// Config represents the main configuration structure
type Config struct {
	Settings      Settings           `yaml:"settings"`
	Services      []Service          `yaml:"services"`
	Notifications NotificationConfig `yaml:"notifications,omitempty"`
}

// LoadConfig reads the configuration file and returns a Config struct
func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// Set default values if not specified
	if config.Settings.CheckInterval == 0 {
		config.Settings.CheckInterval = 1 * time.Minute
	}
	
	if config.Settings.DefaultTimeout == 0 {
		config.Settings.DefaultTimeout = 5 * time.Second
	}

	// Set default storage path if not specified
	if config.Settings.StoragePath == "" {
		// Use the same directory as the config file
		configDir := filepath.Dir(filePath)
		config.Settings.StoragePath = filepath.Join(configDir, "sentinel_status.json")
	}

	return &config, nil
}

// GetServiceTimeout returns the timeout for a specific service
// If the service has a custom timeout, it returns that value
// Otherwise, it returns the default timeout from settings
func (c *Config) GetServiceTimeout(service Service) time.Duration {
	if service.Timeout > 0 {
		return service.Timeout
	}
	return c.Settings.DefaultTimeout
}

// GetNotifiers returns a list of configured notifiers
func (c *Config) GetNotifiers() []notifier.Notifier {
	var notifiers []notifier.Notifier

	// Add Discord notifier if configured
	if c.Notifications.Discord != nil && c.Notifications.Discord.URL != "" {
		notifiers = append(notifiers, notifier.NewDiscordNotifier(*c.Notifications.Discord))
	}

	// Add Slack notifier if configured
	if c.Notifications.Slack != nil && c.Notifications.Slack.URL != "" {
		notifiers = append(notifiers, notifier.NewSlackNotifier(*c.Notifications.Slack))
	}

	// Add generic webhook notifier if configured
	if c.Notifications.Webhook != nil && c.Notifications.Webhook.URL != "" {
		notifiers = append(notifiers, notifier.NewWebhookNotifier(*c.Notifications.Webhook))
	}

	return notifiers
}