package common

import (
	"os"
	"os/exec"
	"regexp"

	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

// GetDBPath tries to determine the database path
// It first tries to get it from the config, then from atuin info
// If no path can be determined, it returns an empty string
func GetDBPath() (string, error) {
	// Try to get the path from atuin info
	dbPath, err := getDBPathFromAtuinInfo()
	if err == nil && dbPath != "" {
		return dbPath, nil
	}

	// If we get here, we couldn't determine the path
	// The caller should handle this by prompting the user
	return "", nil
}

// getDBPathFromAtuinInfo tries to get the database path from the atuin info command
func getDBPathFromAtuinInfo() (string, error) {
	cmd := exec.Command("atuin", "info")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Parse the output to find the client db path
	outputStr := string(output)
	re := regexp.MustCompile(`client db path: "([^"]+)"`)
	matches := re.FindStringSubmatch(outputStr)

	if len(matches) >= 2 {
		dbPath := matches[1]
		// Verify the path exists
		if _, err := os.Stat(dbPath); err == nil {
			return dbPath, nil
		}
	}

	return "", nil
}
