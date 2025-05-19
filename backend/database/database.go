package database

import (
	"backend/common"
	"database/sql"
	"errors"
	"log"
	"os"

	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

var (
	// DB is the global database connection
	DB *sql.DB

	// ErrNoDatabasePath is returned when no database path could be determined
	ErrNoDatabasePath = errors.New("no database path could be determined")

	// ErrDatabaseNotFound is returned when the database file doesn't exist
	ErrDatabaseNotFound = errors.New("database file not found")
)

// Initialize sets up the database connection and ensures tables exist
func Initialize() error {
	// Try to get the database path
	dbPath, err := common.GetDBPath()
	if err != nil {
		return err
	}

	// If we couldn't determine the path, return an error
	if dbPath == "" {
		log.Println("No database path found.")
		return ErrNoDatabasePath
	}

	// Check if the database file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Printf("Database file not found at: %s", dbPath)
		return ErrDatabaseNotFound
	}

	log.Printf("Using database at: %s", dbPath)

	// Connect to SQLite database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return err
	}

	DB = db
	log.Println("Database connection established")

	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// Reinitialize closes the current database connection and initializes a new one
func Reinitialize() error {
	// Close the current connection if it exists
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("Warning: Failed to close database connection: %v", err)
			// Continue anyway, as we want to reinitialize
		}
		DB = nil
	}

	// Initialize a new connection
	return Initialize()
}
