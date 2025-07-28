package checker

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestServiceStatusString(t *testing.T) {
	// Test UP status
	upStatus := ServiceStatus{
		Name:         "TestService",
		URL:          "https://example.com",
		IsUp:         true,
		ResponseTime: 150 * time.Millisecond,
		StatusCode:   200,
	}

	upString := upStatus.String()
	if upString == "" {
		t.Error("Expected non-empty string for UP status")
	}

	if !strings.Contains(upString, "UP") {
		t.Errorf("Expected UP status string to contain 'UP', got: %s", upString)
	}

	if !strings.Contains(upString, "150") {
		t.Errorf("Expected UP status string to contain response time '150', got: %s", upString)
	}

	// Test DOWN status with error
	downStatus := ServiceStatus{
		Name:         "TestService",
		URL:          "https://example.com",
		IsUp:         false,
		ResponseTime: 0,
		StatusCode:   0,
		Error:        &testError{"connection failed"},
	}

	downString := downStatus.String()
	if downString == "" {
		t.Error("Expected non-empty string for DOWN status")
	}

	if !strings.Contains(downString, "DOWN") {
		t.Errorf("Expected DOWN status string to contain 'DOWN', got: %s", downString)
	}

	if !strings.Contains(downString, "connection failed") {
		t.Errorf("Expected DOWN status string to contain error message, got: %s", downString)
	}
}

func TestCheckService(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	// Test successful service check
	status := CheckService("TestService", server.URL)

	if !status.IsUp {
		t.Error("Expected service to be UP")
	}

	if status.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", status.StatusCode)
	}

	if status.Error != nil {
		t.Errorf("Expected no error, got %v", status.Error)
	}

	// Test service check with invalid URL
	badStatus := CheckService("BadService", "http://invalid-url-that-does-not-exist.example")

	if badStatus.IsUp {
		t.Error("Expected service to be DOWN")
	}

	if badStatus.Error == nil {
		t.Error("Expected an error, got nil")
	}
}

// Simple error implementation for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
