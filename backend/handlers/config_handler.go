package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite" // Pure Go SQLite driver

	"backend/common"
	"backend/database"
)

// CheckDatabaseStatus checks if the database is available and working
func CheckDatabaseStatus(c *gin.Context) {
	// First check if we have a database path
	path, err := common.GetDBPath()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to get database path",
			"status": "error",
		})
		return
	}

	// Check if the path exists
	if path == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  "No database path configured",
			"status": "no_path",
		})
		return
	}

	// Check if the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  "Database file not found",
			"status": "not_found",
			"path":   path,
		})
		return
	}

	// Check if the database is accessible
	if database.DB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":  "Database connection not initialized",
			"status": "not_connected",
		})
		return
	}

	// Try to ping the database
	if err := database.DB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":  "Failed to connect to database: " + err.Error(),
			"status": "connection_error",
		})
		return
	}

	// All checks passed
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"path":    path,
		"message": "Database is available and working",
	})
}

// verifyAtuinDatabase checks if the database has the expected Atuin history table structure
func verifyAtuinDatabase(db *sql.DB) error {
	// Check if the history table exists
	var tableName string
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='history'").Scan(&tableName)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("this database does not contain an Atuin history table")
		}
		return fmt.Errorf("failed to verify database structure: %v", err)
	}

	// Verify the table has the expected columns
	// This is a basic check, you might want to make it more thorough
	expectedColumns := []string{"id", "timestamp", "duration", "exit", "command", "cwd", "session", "hostname"}
	for _, column := range expectedColumns {
		query := fmt.Sprintf("SELECT %s FROM history LIMIT 1", column)
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("the history table is missing the expected column '%s'", column)
		}
	}

	return nil
}
