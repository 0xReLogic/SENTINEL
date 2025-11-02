package storage

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/0xReLogic/SENTINEL/checker"
)

const (
	testDBPath           = ":memory:"
	testServiceName      = "Test Service"
	testServiceURL       = "https://example.com"
	errMsgCreateStorage  = "Failed to create storage: %v"
	errMsgSaveCheck      = "Failed to save check: %v"
	errMsgGetHistory     = "Failed to get history: %v"
)

func TestNewSQLiteStorage(t *testing.T) {
	dbPath := testDBPath
	store, err := NewSQLiteStorage(dbPath)
	if err != nil {
		t.Fatalf("Failed to create SQLite storage: %v", err)
	}
	defer store.Close()

	if store.db == nil {
		t.Error("Expected database connection to be initialized")
	}
}

func TestSaveCheck(t *testing.T) {
	store, err := NewSQLiteStorage(testDBPath)
	if err != nil {
		t.Fatalf(errMsgCreateStorage, err)
	}
	defer store.Close()

	check := checker.ServiceStatus{
		Name:         testServiceName,
		URL:          testServiceURL,
		IsUp:         true,
		ResponseTime: 100 * time.Millisecond,
		StatusCode:   200,
		Error:        nil,
	}

	err = store.SaveCheck(check)
	if err != nil {
		t.Errorf(errMsgSaveCheck, err)
	}
}

func TestGetHistory(t *testing.T) {
	store, err := NewSQLiteStorage(testDBPath)
	if err != nil {
		t.Fatalf(errMsgCreateStorage, err)
	}
	defer store.Close()

	// Save multiple checks
	for i := 0; i < 5; i++ {
		check := checker.ServiceStatus{
			Name:         testServiceName,
			URL:          testServiceURL,
			IsUp:         i%2 == 0,
			ResponseTime: time.Duration(i*100) * time.Millisecond,
			StatusCode:   200,
		}
		if err := store.SaveCheck(check); err != nil {
			t.Fatalf(errMsgSaveCheck, err)
		}
	}

	// Retrieve history
	records, err := store.GetHistory(testServiceName, 10)
	if err != nil {
		t.Fatalf(errMsgGetHistory, err)
	}

	if len(records) != 5 {
		t.Errorf("Expected 5 records, got %d", len(records))
	}

	// Verify records are in descending order (newest first)
	for i := 0; i < len(records)-1; i++ {
		if records[i].CheckedAt.Before(records[i+1].CheckedAt) {
			t.Error("Records should be in descending order by checked_at")
		}
	}
}

func TestGetHistoryLimit(t *testing.T) {
	store, err := NewSQLiteStorage(testDBPath)
	if err != nil {
		t.Fatalf(errMsgCreateStorage, err)
	}
	defer store.Close()

	// Save 10 checks
	for i := 0; i < 10; i++ {
		check := checker.ServiceStatus{
			Name:         testServiceName,
			URL:          testServiceURL,
			IsUp:         true,
			ResponseTime: 100 * time.Millisecond,
			StatusCode:   200,
		}
		if err := store.SaveCheck(check); err != nil {
			t.Fatalf(errMsgSaveCheck, err)
		}
	}

	// Retrieve only 5
	records, err := store.GetHistory(testServiceName, 5)
	if err != nil {
		t.Fatalf(errMsgGetHistory, err)
	}

	if len(records) != 5 {
		t.Errorf("Expected 5 records (limit), got %d", len(records))
	}
}

func TestGetHistoryNonExistent(t *testing.T) {
	store, err := NewSQLiteStorage(testDBPath)
	if err != nil {
		t.Fatalf(errMsgCreateStorage, err)
	}
	defer store.Close()

	records, err := store.GetHistory("NonExistent Service", 10)
	if err != nil {
		t.Fatalf(errMsgGetHistory, err)
	}

	if len(records) != 0 {
		t.Errorf("Expected 0 records for non-existent service, got %d", len(records))
	}
}

func TestCleanup(t *testing.T) {
	dbPath := "./test_cleanup.db"
	defer os.Remove(dbPath)

	store, err := NewSQLiteStorage(dbPath)
	if err != nil {
		t.Fatalf(errMsgCreateStorage, err)
	}
	defer store.Close()

	// Insert old record (simulate by direct SQL)
	oldDate := time.Now().AddDate(0, 0, -35) // 35 days ago
	_, err = store.db.Exec(`
		INSERT INTO checks (service_name, service_url, is_up, status_code, response_time_ms, checked_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, "Old Service", "https://old.com", true, 200, 100, oldDate)
	if err != nil {
		t.Fatalf("Failed to insert old record: %v", err)
	}

	// Insert recent record
	check := checker.ServiceStatus{
		Name:         "Recent Service",
		URL:          "https://recent.com",
		IsUp:         true,
		ResponseTime: 100 * time.Millisecond,
		StatusCode:   200,
	}
	if err := store.SaveCheck(check); err != nil {
		t.Fatalf("Failed to save recent check: %v", err)
	}

	// Cleanup records older than 30 days
	err = store.Cleanup(30)
	if err != nil {
		t.Errorf("Cleanup failed: %v", err)
	}

	// Verify old record is deleted
	oldRecords, _ := store.GetHistory("Old Service", 10)
	if len(oldRecords) != 0 {
		t.Errorf("Expected old records to be deleted, got %d", len(oldRecords))
	}

	// Verify recent record still exists
	recentRecords, _ := store.GetHistory("Recent Service", 10)
	if len(recentRecords) != 1 {
		t.Errorf("Expected recent record to remain, got %d", len(recentRecords))
	}
}

func TestSaveCheckWithError(t *testing.T) {
	store, err := NewSQLiteStorage(testDBPath)
	if err != nil {
		t.Fatalf(errMsgCreateStorage, err)
	}
	defer store.Close()

	check := checker.ServiceStatus{
		Name:         "Failed Service",
		URL:          "https://failed.com",
		IsUp:         false,
		ResponseTime: 0,
		StatusCode:   0,
		Error:        fmt.Errorf("connection timeout"),
	}

	err = store.SaveCheck(check)
	if err != nil {
		t.Errorf("Failed to save check with error: %v", err)
	}

	records, _ := store.GetHistory("Failed Service", 1)
	if len(records) == 0 {
		t.Fatal("Expected to retrieve saved record")
	}

	if records[0].ErrorMessage != "connection timeout" {
		t.Errorf("Expected error message 'connection timeout', got '%s'", records[0].ErrorMessage)
	}
}
