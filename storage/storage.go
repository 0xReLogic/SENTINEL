package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sentinel/checker"
)

// StatusStorage is responsible for storing and retrieving service statuses
type StatusStorage struct {
	statuses map[string]checker.ServiceStatus
	mutex    sync.RWMutex
	filePath string
}

// NewStatusStorage creates a new StatusStorage
func NewStatusStorage(storagePath string) (*StatusStorage, error) {
	// Create storage directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(storagePath), 0755); err != nil {
		return nil, fmt.Errorf("error creating storage directory: %w", err)
	}

	storage := &StatusStorage{
		statuses: make(map[string]checker.ServiceStatus),
		filePath: storagePath,
	}

	// Try to load existing statuses
	if err := storage.Load(); err != nil {
		// If the file doesn't exist, that's fine, we'll create it later
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error loading statuses: %w", err)
		}
	}

	return storage, nil
}

// GetStatus retrieves the status for a service
func (s *StatusStorage) GetStatus(serviceName string) (checker.ServiceStatus, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	status, exists := s.statuses[serviceName]
	return status, exists
}

// SetStatus stores the status for a service and returns the old status if it exists
func (s *StatusStorage) SetStatus(status checker.ServiceStatus) (checker.ServiceStatus, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	oldStatus, exists := s.statuses[status.Name]
	s.statuses[status.Name] = status

	// Save to disk after updating
	if err := s.save(); err != nil {
		fmt.Printf("Warning: Failed to save statuses: %v\n", err)
	}

	return oldStatus, exists
}

// HasStatusChanged checks if the status has changed from the previous status
func (s *StatusStorage) HasStatusChanged(serviceName string, newStatus checker.ServiceStatus) (bool, checker.ServiceStatus) {
	oldStatus, exists := s.GetStatus(serviceName)
	if !exists {
		return false, oldStatus
	}

	return oldStatus.IsUp != newStatus.IsUp, oldStatus
}

// Load loads statuses from disk
func (s *StatusStorage) Load() error {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	return json.Unmarshal(data, &s.statuses)
}

// save saves statuses to disk
func (s *StatusStorage) save() error {
	data, err := json.MarshalIndent(s.statuses, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling statuses: %w", err)
	}

	return os.WriteFile(s.filePath, data, 0644)
}

// Cleanup removes old statuses that haven't been updated in a while
func (s *StatusStorage) Cleanup(maxAge time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	for name, status := range s.statuses {
		// If the status has a timestamp and it's older than maxAge, remove it
		if !status.Timestamp.IsZero() && now.Sub(status.Timestamp) > maxAge {
			delete(s.statuses, name)
		}
	}

	// Save changes
	if err := s.save(); err != nil {
		fmt.Printf("Warning: Failed to save statuses during cleanup: %v\n", err)
	}
}
