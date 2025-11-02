package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/0xReLogic/SENTINEL/checker"
	_ "modernc.org/sqlite"
)

// SQLiteStorage implements Storage interface using SQLite
type SQLiteStorage struct {
	db *sql.DB
}

const schema = `
CREATE TABLE IF NOT EXISTS checks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    service_name TEXT NOT NULL,
    service_url TEXT NOT NULL,
    is_up BOOLEAN NOT NULL,
    status_code INTEGER,
    response_time_ms INTEGER,
    error_message TEXT,
    checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_service_time ON checks(service_name, checked_at);
CREATE INDEX IF NOT EXISTS idx_checked_at ON checks(checked_at);
`

// NewSQLiteStorage creates a new SQLite storage instance
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Enable WAL mode for better concurrency
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Create tables and indexes
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return &SQLiteStorage{db: db}, nil
}

// SaveCheck saves a service check result to the database
func (s *SQLiteStorage) SaveCheck(check checker.ServiceStatus) error {
	var errorMsg string
	if check.Error != nil {
		errorMsg = check.Error.Error()
	}

	query := `
		INSERT INTO checks (service_name, service_url, is_up, status_code, response_time_ms, error_message)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(query,
		check.Name,
		check.URL,
		check.IsUp,
		check.StatusCode,
		check.ResponseTime.Milliseconds(),
		errorMsg,
	)

	if err != nil {
		return fmt.Errorf("failed to save check: %w", err)
	}

	return nil
}

// GetHistory retrieves check history for a service
func (s *SQLiteStorage) GetHistory(serviceName string, limit int) ([]CheckRecord, error) {
	query := `
		SELECT id, service_name, service_url, is_up, status_code, response_time_ms, error_message, checked_at
		FROM checks
		WHERE service_name = ?
		ORDER BY checked_at DESC
		LIMIT ?
	`

	rows, err := s.db.Query(query, serviceName, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query history: %w", err)
	}
	defer rows.Close()

	var records []CheckRecord
	for rows.Next() {
		var r CheckRecord
		err := rows.Scan(
			&r.ID,
			&r.ServiceName,
			&r.ServiceURL,
			&r.IsUp,
			&r.StatusCode,
			&r.ResponseTimeMs,
			&r.ErrorMessage,
			&r.CheckedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		records = append(records, r)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return records, nil
}

// Cleanup removes old records based on retention policy
func (s *SQLiteStorage) Cleanup(retentionDays int) error {
	query := `DELETE FROM checks WHERE checked_at < datetime('now', '-' || ? || ' days')`

	result, err := s.db.Exec(query, retentionDays)
	if err != nil {
		return fmt.Errorf("failed to cleanup old records: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		fmt.Printf("Cleaned up %d old records (older than %d days)\n", rowsAffected, retentionDays)
	}

	return nil
}

// Close closes the database connection
func (s *SQLiteStorage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
