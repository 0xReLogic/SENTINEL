package config

import (
	// "fmt"
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

// config/config_test.go

// func TestLoadConfigWithEnvVars(t *testing.T) {
//     // 1. Define Test Environment Variables
//     const testToken = "TEST_BOT_TOKEN_12345"
//     const testChatID = "123456789"

//     // Set the environment variables, and use t.Cleanup to ensure they are restored
//     // to their original state after the test finishes. This is critical for isolated testing.
//     t.Setenv("TELEGRAM_BOT_TOKEN", testToken) //
//     t.Setenv("TELEGRAM_CHAT_ID", testChatID) //

//     // 2. Create a temporary config file with placeholders
//     tempDir := t.TempDir()
//     configPath := filepath.Join(tempDir, "env-config.yaml")

//     configContent := fmt.Sprintf(`
// services:
//   - name: "Placeholder Service"
//     url: "https://service-a.com"
// notifications:
//   telegram:
//     enabled: true
//     bot_token: "${TELEGRAM_BOT_TOKEN}"
//     chat_id: "%s" # Test interpolation with and without brackets
//     notify_on: ["down", "recovery"]
// `, "${TELEGRAM_CHAT_ID}")

//     err := os.WriteFile(configPath, []byte(configContent), 0644)
//     if err != nil {
//         t.Fatalf("Failed to write test config: %v", err)
//     }

//     // 3. Load the config (which should now call os.ExpandEnv inside LoadConfig)
//     config, err := LoadConfig(configPath)
//     if err != nil {
//         t.Fatalf("LoadConfig failed with environment variables: %v", err)
//     }

//     // 4. Verification (The main assertion)
//     // Check if the service data is intact
//     if len(config.Services) != 1 {
//         t.Errorf("Expected 1 service, got %d", len(config.Services))
//     }

//     // Check if the environment variables were correctly loaded and replaced
//     if config.Notifications.Telegram.BotToken != testToken {
//         t.Errorf("BotToken interpolation failed. Expected '%s', got '%s'",
//             testToken, config.Notifications.Telegram.BotToken)
//     }

//     if config.Notifications.Telegram.ChatID != testChatID {
//         t.Errorf("ChatID interpolation failed. Expected '%s', got '%s'",
//             testChatID, config.Notifications.Telegram.ChatID)
//     }

//     // Check if the other fields are loaded correctly
//     if !config.Notifications.Telegram.Enabled {
//         t.Error("Expected Telegram to be enabled, got false")
//     }
//     if config.Notifications.Telegram.NotifyOn[0] != "down" {
//         t.Error("Expected NotifyOn to include 'down'")
//     }
// }
