package models

import (
	"time"
)

// History represents a record in the history table
type History struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Duration  int    `json:"duration"`
	Exit      int    `json:"exit"`
	Command   string `json:"command"`
	Cwd       string `json:"cwd"`
	Session   string `json:"session"`
	Hostname  string `json:"hostname"`
	DeletedAt *int64 `json:"deleted_at,omitempty"`
}

// HistoryFilter represents filter parameters for history queries
type HistoryFilter struct {
	ID        string `form:"id"`
	Command   string `form:"command"`
	Cwd       string `form:"cwd"`
	Hostname  string `form:"hostname"`
	StartTime int64  `form:"start_time"`
	EndTime   int64  `form:"end_time"`
	Exit      *int   `form:"exit"`
	Limit     int    `form:"limit"`
	Offset    int    `form:"offset"`
}

// FormatTime formats a nanosecond timestamp as a human-readable string
func (h *History) FormatTime() string {
	// Convert nanoseconds to Unix time
	seconds := h.Timestamp / 1000000000
	nanoseconds := h.Timestamp % 1000000000
	t := time.Unix(seconds, nanoseconds)
	return t.Format(time.RFC3339)
}

// FormatDuration formats the duration in a human-readable format
func (h *History) FormatDuration() string {
	// Convert nanoseconds directly to time.Duration (which is in nanoseconds)
	duration := time.Duration(h.Duration / 1000000000)
	return duration.String()
}
