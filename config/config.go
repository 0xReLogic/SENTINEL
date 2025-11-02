// Package config provides functionality for loading and managing configuration
// Repository: https://github.com/0xReLogic/SENTINEL
package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
	
)

// Default configuration values
const (
	DefaultInterval = 1 * time.Minute
	DefaultTimeout  = 5 * time.Second
)



// Service represents a single service to be monitored
type Service struct {
	Name     string        `yaml:"name"`
	URL      string        `yaml:"url"`
	Interval time.Duration `yaml:"interval"`
	Timeout  time.Duration `yaml:"timeout"`
}

// Config represents the main configuration structure
type Config struct {
	Services      []Service          `yaml:"services"`
	Notifications NotificationConfig `yaml:"notifications"`
	Storage       StorageConfig      `yaml:"storage"`
}

type TelegramConfig struct {
	Enabled  bool     `yaml:"enabled"`
	BotToken string   `yaml:"bot_token"`
	ChatID   string   `yaml:"chat_id"`
	NotifyOn []string `yaml:"notify_on"`
}

type DiscordConfig struct {
	Enabled    bool     `yaml:"enabled"`
	WebhookURL string   `yaml:"webhook_url"`
	NotifyOn   []string `yaml:"notify_on"`
}

type NotificationConfig struct {
	Telegram TelegramConfig `yaml:"telegram"`
	Discord  DiscordConfig  `yaml:"discord"`
}

type StorageConfig struct {
	Type          string `yaml:"type"`
	Path          string `yaml:"path"`
	RetentionDays int    `yaml:"retention_days"`
}
// LoadConfig reads the configuration file from the given path, expands any
// environment variables, and unmarshals it into a Config struct.
func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	content := string(data)
	content = os.ExpandEnv(content)
	data = []byte(content)
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// apply defaults for optional fields
	for i := range config.Services {
		svc := &config.Services[i]

		if svc.Interval == 0 {
			svc.Interval = DefaultInterval
		}
		if svc.Timeout == 0 {
			svc.Timeout = DefaultTimeout
		}

		// Validate
		if svc.Interval < 0 {
			return nil, fmt.Errorf("service '%s': interval must be positive, got %v", svc.Name, svc.Interval)
		}
		if svc.Timeout < 0 {
			return nil, fmt.Errorf("service '%s': timeout must be positive, got %v", svc.Name, svc.Timeout)
		}
	}

	return &config, nil
}
