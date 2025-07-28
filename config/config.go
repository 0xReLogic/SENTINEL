package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Service represents a single service to be monitored
type Service struct {
	Name    string        `yaml:"name"`
	URL     string        `yaml:"url"`
	Timeout time.Duration `yaml:"timeout,omitempty"` // Optional timeout per service
}

// Settings represents global application settings
type Settings struct {
	CheckInterval  time.Duration `yaml:"check_interval,omitempty"`  // Interval between checks
	DefaultTimeout time.Duration `yaml:"default_timeout,omitempty"` // Default timeout for all services
}

// Config represents the main configuration structure
type Config struct {
	Settings Settings  `yaml:"settings"`
	Services []Service `yaml:"services"`
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
