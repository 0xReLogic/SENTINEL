package checker

import (
	"fmt"
	"net/http"
	"time"
)

// ServiceStatus represents the result of a service check
type ServiceStatus struct {
	Name         string
	URL          string
	IsUp         bool
	ResponseTime time.Duration
	StatusCode   int
	Error        error
}

// String returns a formatted string representation of the service status
func (s ServiceStatus) String() string {
	status := "UP"
	if !s.IsUp {
		status = "DOWN"
	}

	if s.Error != nil {
		return fmt.Sprintf("[%s] %s - Error: %s", status, s.Name, s.Error)
	}

	return fmt.Sprintf("[%s] %s - %d ms (HTTP %d)", status, s.Name, s.ResponseTime.Milliseconds(), s.StatusCode)
}

// CheckService performs an HTTP GET request to the given URL and returns the service status
func CheckService(name, url string, timeout time.Duration) ServiceStatus {
	result := ServiceStatus{
		Name: name,
		URL:  url,
	}

	// Create HTTP client with configured timeout
	client := &http.Client{
		Timeout: timeout,
	}

	// Record start time
	startTime := time.Now()

	// Send HTTP GET request
	resp, err := client.Get(url)

	// Calculate response time
	result.ResponseTime = time.Since(startTime)

	// Handle errors
	if err != nil {
		result.IsUp = false
		result.Error = err
		return result
	}
	defer resp.Body.Close()

	// Set status code
	result.StatusCode = resp.StatusCode

	// Determine if service is up (2xx or 3xx status codes)
	result.IsUp = resp.StatusCode >= 200 && resp.StatusCode < 400

	return result
}
