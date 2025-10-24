package cmd

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/0xReLogic/SENTINEL/config"
	"github.com/spf13/cobra"
)

func assertNoPanic(t *testing.T) {
	t.Helper()
	if r := recover(); r != nil {
		t.Errorf("%s panicked: %v", t.Name(), r)
	}
}

// createMockServer creates an HTTP test server for testing
func createMockServer(t *testing.T) (*httptest.Server, func()) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch {
		case strings.Contains(r.URL.Path, "/status/200"):
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		case strings.Contains(r.URL.Path, "/status/201"):
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("Created"))
		case strings.Contains(r.URL.Path, "/delay/"):
			time.Sleep(6 * time.Second) // Client timeout is 5 seconds
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Delayed Response"))
		case strings.Contains(r.URL.Path, "/redirect/"):
			w.Header().Set("Location", "/status/200")
			w.WriteHeader(http.StatusMovedPermanently)
		case strings.Contains(r.URL.Path, "/error"):
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		default:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		}
	}))

	// return server and cleanup function
	cleanup := func() {
		server.Close()
	}

	return server, cleanup
}

func createTempConfig(t *testing.T, content string) string {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.yaml")

	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp config: %v", err)
	}

	return configPath
}

func TestConstants(t *testing.T) {

	if appName == "" {
		t.Error("appName constant should not be empty")
	}
	if appRepository == "" {
		t.Error("appRepository constant should not be empty")
	}
	if defaultConfigFile == "" {
		t.Error("defaultConfigFile constant should not be empty")
	}
	if timestampFormat == "" {
		t.Error("timestampFormat constant should not be empty")
	}
	if exitSuccess != 0 {
		t.Errorf("exitSuccess should be 0, got %d", exitSuccess)
	}
	if exitError != 1 {
		t.Errorf("exitError should be 1, got %d", exitError)
	}
	if exitConfigError != 2 {
		t.Errorf("exitConfigError should be 2, got %d", exitConfigError)
	}
}

func TestExecute(t *testing.T) {
	defer assertNoPanic(t)

	if rootCmd == nil {
		t.Fatal("rootCmd should not be nil")
	}

	if rootCmd.Use != appName {
		t.Errorf("Expected rootCmd.Use to be %s, got %s", appName, rootCmd.Use)
	}

	commands := rootCmd.Commands()
	if len(commands) == 0 {
		t.Error("Expected rootCmd to have subcommands registered")
	}

	expectedCommands := []string{cmdNameRun, cmdNameOnce, cmdNameValidate}
	commandNames := make(map[string]bool)
	for _, cmd := range commands {
		commandNames[cmd.Use] = true
	}

	for _, expected := range expectedCommands {
		if !commandNames[expected] {
			t.Errorf("Expected command %s to be registered", expected)
		}
	}

	var runCmd *cobra.Command
	for _, cmd := range commands {
		if cmd.Use == cmdNameRun {
			runCmd = cmd
			break
		}
	}

	if runCmd == nil {
		t.Fatal("Expected run command to be registered")
	}

	if runCmd.Use != cmdNameRun {
		t.Errorf("Expected run command use to be %s, got %s", cmdNameRun, runCmd.Use)
	}
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name           string
		configContent  string
		setupPath      func(t *testing.T) string
		expectError    bool
		validateConfig func(t *testing.T, cfg *config.Config)
	}{
		{
			name: "valid config with two services",
			configContent: `
services:
  - name: "Test Service"
    url: "https://example.com"
  - name: "Another Service"
    url: "https://test.org"
`,
			setupPath: func(t *testing.T) string {
				return createTempConfig(t, `
services:
  - name: "Test Service"
    url: "https://example.com"
  - name: "Another Service"
    url: "https://test.org"
`)
			},
			expectError: false,
			validateConfig: func(t *testing.T, cfg *config.Config) {
				if cfg == nil {
					t.Error("Expected config to be loaded")
					return
				}
				if len(cfg.Services) != 2 {
					t.Errorf("Expected 2 services, got %d", len(cfg.Services))
				}
			},
		},
		{
			name:          "empty config",
			configContent: `services: []`,
			setupPath: func(t *testing.T) string {
				return createTempConfig(t, `services: []`)
			},
			expectError: false,
			validateConfig: func(t *testing.T, cfg *config.Config) {
				if len(cfg.Services) != 0 {
					t.Errorf("Expected 0 services for empty config, got %d", len(cfg.Services))
				}
			},
		},
		{
			name: "invalid YAML",
			configContent: `
services:
  - name: "Test Service"
    url: https:
    invalid_yaml:
      - [
`,
			setupPath: func(t *testing.T) string {
				return createTempConfig(t, `
services:
  - name: "Test Service"
    url: https:
    invalid_yaml:
      - [
`)
			},
			expectError: true,
			validateConfig: func(t *testing.T, cfg *config.Config) {
			},
		},
		{
			name:          "non-existent file",
			configContent: "",
			setupPath: func(t *testing.T) string {
				return "/non/existent/file.yaml"
			},
			expectError: true,
			validateConfig: func(t *testing.T, cfg *config.Config) {
			},
		},
		{
			name:          "weird path",
			configContent: "",
			setupPath: func(t *testing.T) string {
				return "./././sentinel.yaml"
			},
			expectError: true,
			validateConfig: func(t *testing.T, cfg *config.Config) {
			},
		},
		{
			name: "comment-only config",
			configContent: `
# This is just a comment
services: []
`,
			setupPath: func(t *testing.T) string {
				return createTempConfig(t, `
# This is just a comment
services: []
`)
			},
			expectError: false,
			validateConfig: func(t *testing.T, cfg *config.Config) {
				if len(cfg.Services) != 0 {
					t.Errorf("Expected 0 services for comment config, got %d", len(cfg.Services))
				}
			},
		},
		{
			name: "special characters in service names and URLs",
			configContent: `
services:
  - name: "Test-Service_123"
    url: "https://example.com/path?query=value&other=test"
  - name: "Service with spaces"
    url: "https://test.com:8080/api/v1"
`,
			setupPath: func(t *testing.T) string {
				return createTempConfig(t, `
services:
  - name: "Test-Service_123"
    url: "https://example.com/path?query=value&other=test"
  - name: "Service with spaces"
    url: "https://test.com:8080/api/v1"
`)
			},
			expectError: false,
			validateConfig: func(t *testing.T, cfg *config.Config) {
				if len(cfg.Services) != 2 {
					t.Errorf("Expected 2 services, got %d", len(cfg.Services))
					return
				}
				if cfg.Services[0].Name != "Test-Service_123" {
					t.Errorf("Expected first service name 'Test-Service_123', got '%s'", cfg.Services[0].Name)
				}
				if cfg.Services[1].URL != "https://test.com:8080/api/v1" {
					t.Errorf("Expected second service URL 'https://test.com:8080/api/v1', got '%s'", cfg.Services[1].URL)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setupPath(t)

			cfg, err := loadConfig(path)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("loadConfig failed: %v", err)
				}
				tt.validateConfig(t, cfg)
			}
		})
	}
}

func TestValidateServices(t *testing.T) {
	tests := []struct {
		name     string
		services []config.Service
		wantErr  bool
	}{
		{
			name: "valid services",
			services: []config.Service{
				{Name: "Test", URL: "https://example.com", Interval: config.DefaultInterval, Timeout: config.DefaultTimeout},
			},
			wantErr: false,
		},
		{
			name: "empty name",
			services: []config.Service{
				{Name: "", URL: "https://example.com", Interval: config.DefaultInterval, Timeout: config.DefaultTimeout},
			},
			wantErr: true,
		},
		{
			name: "empty URL",
			services: []config.Service{
				{Name: "Test", URL: "", Interval: config.DefaultInterval, Timeout: config.DefaultTimeout},
			},
			wantErr: true,
		},
		{
			name: "invalid URL scheme",
			services: []config.Service{
				{Name: "Test", URL: "ftp://example.com", Interval: config.DefaultInterval, Timeout: config.DefaultTimeout},
			},
			wantErr: true,
		},
		{
			name: "non-positive interval",
			services: []config.Service{
				{Name: "Test", URL: "https://example.com", Interval: 0, Timeout: config.DefaultTimeout},
			},
			wantErr: true,
		},
		{
			name: "non-positive timeout",
			services: []config.Service{
				{Name: "Test", URL: "https://example.com", Interval: config.DefaultInterval, Timeout: 0},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validateServices(tt.services)
			gotErr := len(errors) > 0
			if gotErr != tt.wantErr {
				t.Errorf("validateServices() error = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		url   string
		valid bool
	}{
		{"https://example.com", true},
		{"http://example.com", true},
		{"https://example.com:8080", true},
		{"https://example.com/path", true},
		{"ftp://example.com", false},
		{"example.com", false},
		{"", false},
		{"https://", false},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			result := isValidURL(tt.url)
			if result != tt.valid {
				t.Errorf("isValidURL(%q) = %v, want %v", tt.url, result, tt.valid)
			}
		})
	}
}

func TestPrintBanner(t *testing.T) {
	defer assertNoPanic(t)

	cfg := &config.Config{
		Services: []config.Service{
			{Name: "Test", URL: "https://example.com"},
		},
	}

	printBanner(cfg)
}

func TestRunChecks(t *testing.T) {
	defer assertNoPanic(t)

	// mock HTTP server
	server, cleanup := createMockServer(t)
	defer cleanup()

	cfg := &config.Config{
		Services: []config.Service{
			{Name: "Test", URL: server.URL + "/status/200"},
		},
	}
	stateManager := NewStateManager() 
	runChecksAndGetStatus(cfg, stateManager)
}

func TestRunChecksAndGetStatus(t *testing.T) {
	defer assertNoPanic(t)

	// mock HTTP server
	server, cleanup := createMockServer(t)
	defer cleanup()

	tests := []struct {
		name     string
		services []config.Service
		expected bool
	}{
		{
			name: "single service - UP",
			services: []config.Service{
				{Name: "Test", URL: server.URL + "/status/200"},
			},
			expected: true,
		},
		{
			name: "multiple services - all UP",
			services: []config.Service{
				{Name: "Test1", URL: server.URL + "/status/200"},
				{Name: "Test2", URL: server.URL + "/status/201"},
			},
			expected: true,
		},
		{
			name: "service with timeout - DOWN",
			services: []config.Service{
				{Name: "TimeoutTest", URL: server.URL + "/delay/10"},
			},
			expected: false,
		},
		{
			name: "service with redirect - UP",
			services: []config.Service{
				{Name: "RedirectTest", URL: server.URL + "/redirect/1"},
			},
			expected: true,
		},
		{
			name:     "empty services - UP",
			services: []config.Service{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{Services: tt.services}

			stateManager :=NewStateManager()
			result := runChecksAndGetStatus(cfg, stateManager )

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
