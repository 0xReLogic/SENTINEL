package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "sentinel-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test config file
	configPath := filepath.Join(tempDir, "test-config.yaml")
	configContent := `
services:
  - name: "Test Service"
    url: "https://example.com"
  - name: "Another Service"
    url: "https://example.org"
`
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Test loading the config
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify the config
	if len(config.Services) != 2 {
		t.Errorf("Expected 2 services, got %d", len(config.Services))
	}

	if config.Services[0].Name != "Test Service" {
		t.Errorf("Expected service name 'Test Service', got '%s'", config.Services[0].Name)
	}

	if config.Services[0].URL != "https://example.com" {
		t.Errorf("Expected service URL 'https://example.com', got '%s'", config.Services[0].URL)
	}

	if config.Services[1].Name != "Another Service" {
		t.Errorf("Expected service name 'Another Service', got '%s'", config.Services[1].Name)
	}

	if config.Services[1].URL != "https://example.org" {
		t.Errorf("Expected service URL 'https://example.org', got '%s'", config.Services[1].URL)
	}
}

func TestLoadConfigInvalidPath(t *testing.T) {
	_, err := LoadConfig("non-existent-file.yaml")
	if err == nil {
		t.Error("Expected error when loading non-existent file, got nil")
	}
}

func TestLoadConfigInvalidYAML(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "sentinel-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create an invalid YAML file
	configPath := filepath.Join(tempDir, "invalid-config.yaml")
	invalidContent := `
services:
  - name: "Invalid Service"
    url: https://example.com
    invalid_yaml:
      - [
`
	err = os.WriteFile(configPath, []byte(invalidContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid config: %v", err)
	}

	// Test loading the invalid config
	_, err = LoadConfig(configPath)
	if err == nil {
		t.Error("Expected error when loading invalid YAML, got nil")
	}
}

func TestEmptyConfig(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "sentinel-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create an empty config file
	configPath := filepath.Join(tempDir, "empty-config.yaml")
	emptyContent := `
services: []
`
	err = os.WriteFile(configPath, []byte(emptyContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write empty config: %v", err)
	}

	// Test loading the empty config
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed for empty config: %v", err)
	}

	if len(config.Services) != 0 {
		t.Errorf("Expected 0 services for empty config, got %d", len(config.Services))
	}
}
