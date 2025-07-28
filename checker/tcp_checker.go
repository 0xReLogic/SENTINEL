package checker

import (
	"fmt"
	"net"
	"time"
)

// CheckTCPPort performs a TCP connection check to the given host and port
func CheckTCPPort(name, host string, port int, timeout time.Duration) ServiceStatus {
	result := ServiceStatus{
		Name:      name,
		URL:       fmt.Sprintf("%s:%d", host, port),
		Timestamp: time.Now(),
	}

	// Create address string
	address := fmt.Sprintf("%s:%d", host, port)

	// Record start time
	startTime := time.Now()

	// Try to establish TCP connection
	conn, err := net.DialTimeout("tcp", address, timeout)
	
	// Calculate response time
	result.ResponseTime = time.Since(startTime)

	// Handle errors
	if err != nil {
		result.IsUp = false
		result.Error = err
		return result
	}
	defer conn.Close()

	// If we got here, the connection was successful
	result.IsUp = true

	return result
}