// Package config provides functionality for loading and managing configuration
// Repository: https://github.com/0xReLogic/SENTINEL
package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// Service represents a single service to be monitored
type Service struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

// Config represents the main configuration structure
type Config struct {
	Services      []Service          `yaml:"services"`
	Notifications NotificationConfig `yaml:"notifications"`
}

type TelegramConfig struct {
	Enabled  bool     `yaml:"enabled"`
	BotToken string   `yaml:"bot_token"`
	ChatID   string   `yaml:"chat_id"`
	NotifyOn []string `yaml:"notify_on"`
}

type NotificationConfig struct {
	Telegram TelegramConfig `yaml:"telegram"`
}

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

	return &config, nil
}
