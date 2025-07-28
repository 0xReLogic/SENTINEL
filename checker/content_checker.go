package checker

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// CheckContent performs an HTTP GET request and checks if the response contains the expected string
func CheckContent(name, url, expectString string, timeout time.Duration) ServiceStatus {
	// First perform a regular HTTP check
	result := CheckService(name, url, timeout)
	
	// If the service is down, no need to check content
	if !result.IsUp || result.Error != nil {
		return result
	}
	
	// Create HTTP client with timeout
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
	
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.IsUp = false
		result.Error = fmt.Errorf("error reading response body: %w", err)
		return result
	}
	
	// Check if response contains the expected string
	if !strings.Contains(string(body), expectString) {
		result.IsUp = false
		result.Error = fmt.Errorf("response does not contain expected string: %s", expectString)
	}
	
	return result
}