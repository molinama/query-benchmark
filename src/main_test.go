package main

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func Test_run(t *testing.T) {
	csvContent := `hostname,start_time,end_time
host_000008,2017-01-01 08:59:22,2017-01-01 09:59:22
host_000001,2017-01-02 13:02:02,2017-01-02 14:02:02
host_000008,2017-01-02 18:50:28,2017-01-02 19:50:28`
	tmpCSV, err := os.CreateTemp("", "test-*.csv")
	if err != nil {
		t.Fatalf("Error creating temp CSV file: %v", err)
	}
	defer os.Remove(tmpCSV.Name())

	if _, err := tmpCSV.Write([]byte(csvContent)); err != nil {
		t.Fatalf("Error writing to temp CSV file: %v", err)
	}
	if err := tmpCSV.Close(); err != nil {
		t.Fatalf("Error closing temp CSV file: %v", err)
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE cpu_usage (
		ts TIMESTAMPTZ,
		host TEXT,
		usage DOUBLE PRECISION
	)`)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	config := Config{
		csvFilePath:   tmpCSV.Name(),
		numberWorkers: 20,
		dbConnString:  ":memory:",
		db:            db,
	}

	err = run(config)
	if err != nil {
		t.Fatalf("Error running main function: %v", err)
	}
}
