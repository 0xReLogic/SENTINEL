package storage

import (
	"time"

	"github.com/0xReLogic/SENTINEL/checker"
)

// Storage interface defines methods for persisting check history
type Storage interface {
	// SaveCheck saves a service check result to storage
	SaveCheck(check checker.ServiceStatus) error

	// GetHistory retrieves check history for a service
	GetHistory(serviceName string, limit int) ([]CheckRecord, error)

	// Cleanup removes old records based on retention policy
	Cleanup(retentionDays int) error

	// Close closes the storage connection
	Close() error
}

// CheckRecord represents a stored check result
type CheckRecord struct {
	ID             int64
	ServiceName    string
	ServiceURL     string
	IsUp           bool
	StatusCode     int
	ResponseTimeMs int64
	ErrorMessage   string
	CheckedAt      time.Time
}
