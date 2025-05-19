package repositories

import (
	"fmt"
	"strings"
	"time"

	"backend/database"
	"backend/models"
)

// GetHistoryEntries retrieves history entries based on filter criteria
func GetHistoryEntries(filter models.HistoryFilter) ([]models.History, error) {
	entries := make([]models.History, 0)

	// Build the query
	query := "SELECT id, timestamp, duration, exit, command, cwd, session, hostname, deleted_at FROM history WHERE deleted_at IS NULL"
	params := []interface{}{}

	// Apply filters
	if filter.ID != "" {
		query += " AND id = ?"
		params = append(params, filter.ID)
	}

	if filter.Command != "" {
		query += " AND command LIKE ?"
		params = append(params, "%"+filter.Command+"%")
	}

	if filter.Cwd != "" {
		query += " AND cwd LIKE ?"
		params = append(params, "%"+filter.Cwd+"%")
	}

	if filter.Hostname != "" {
		query += " AND hostname = ?"
		params = append(params, filter.Hostname)
	}

	if filter.StartTime > 0 {
		query += " AND timestamp >= ?"
		params = append(params, filter.StartTime)
	}

	if filter.EndTime > 0 {
		query += " AND timestamp <= ?"
		params = append(params, filter.EndTime)
	}

	if filter.Exit != nil {
		query += " AND exit = ?"
		params = append(params, *filter.Exit)
	}

	// Add order by
	query += " ORDER BY timestamp DESC"

	// Add limit and offset
	if filter.Limit > 0 {
		query += " LIMIT ?"
		params = append(params, filter.Limit)

		if filter.Offset > 0 {
			query += " OFFSET ?"
			params = append(params, filter.Offset)
		}
	}

	// Execute the query
	rows, err := database.DB.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process the results
	for rows.Next() {
		var entry models.History
		var deletedAt *int64

		err := rows.Scan(
			&entry.ID,
			&entry.Timestamp,
			&entry.Duration,
			&entry.Exit,
			&entry.Command,
			&entry.Cwd,
			&entry.Session,
			&entry.Hostname,
			&deletedAt,
		)

		if err != nil {
			return nil, err
		}

		entry.DeletedAt = deletedAt
		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

// BatchDeleteHistoryEntries marks multiple history entries as deleted
func BatchDeleteHistoryEntries(ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	// Create placeholders for the IN clause
	placeholders := make([]string, len(ids))
	params := make([]interface{}, len(ids)+1)

	// Current timestamp for deleted_at in nanoseconds
	now := time.Now().UnixNano()
	params[0] = now

	for i, id := range ids {
		placeholders[i] = "?"
		params[i+1] = id
	}

	// Build the query
	query := fmt.Sprintf(
		"UPDATE history SET deleted_at = ? WHERE id IN (%s)",
		strings.Join(placeholders, ","),
	)

	// Execute the query
	_, err := database.DB.Exec(query, params...)
	return err
}
