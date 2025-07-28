package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Service represents a single service to be monitored
type Service struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

// Config represents the main configuration structure
type Config struct {
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

	return &config, nil
}