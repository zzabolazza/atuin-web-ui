package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/models"
	"backend/repositories"
)

// GetHistoryEntries handles the request to get history entries with filtering
func GetHistoryEntries(c *gin.Context) {
	var filter models.HistoryFilter
	
	// Bind query parameters to filter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter parameters"})
		return
	}
	
	// Set default limit if not provided
	if filter.Limit <= 0 {
		filter.Limit = 100
	}
	
	// Get entries from repository
	entries, err := repositories.GetHistoryEntries(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve history entries"})
		return
	}
	
	c.JSON(http.StatusOK, entries)
}

// BatchDeleteHistoryEntries handles the request to batch delete history entries
func BatchDeleteHistoryEntries(c *gin.Context) {
	var request struct {
		IDs []string `json:"ids" binding:"required"`
	}
	
	// Bind JSON body to request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	
	// Delete entries
	err := repositories.BatchDeleteHistoryEntries(request.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete history entries"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "History entries deleted successfully"})
}
