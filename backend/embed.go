package main

import (
	"embed"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

//go:embed frontend_dist
var frontendFS embed.FS

// ServeFrontend sets up routes to serve the embedded frontend files
func ServeFrontend(router *gin.Engine) {
	// Create a sub-filesystem for the frontend_dist directory
	frontendDistFS, err := fs.Sub(frontendFS, "frontend_dist")
	if err != nil {
		panic(err)
	}

	// Register common MIME types
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".html", "text/html")
	mime.AddExtensionType(".json", "application/json")
	mime.AddExtensionType(".png", "image/png")
	mime.AddExtensionType(".jpg", "image/jpeg")
	mime.AddExtensionType(".jpeg", "image/jpeg")
	mime.AddExtensionType(".svg", "image/svg+xml")
	mime.AddExtensionType(".ico", "image/x-icon")
	mime.AddExtensionType(".woff", "font/woff")
	mime.AddExtensionType(".woff2", "font/woff2")
	mime.AddExtensionType(".ttf", "font/ttf")
	mime.AddExtensionType(".eot", "application/vnd.ms-fontobject")

	// Serve the index.html file for the root path
	router.GET("/", func(c *gin.Context) {
		serveIndexHTML(c, frontendDistFS)
	})

	// Serve static files directly with explicit routes
	router.GET("/css/*filepath", serveStaticFile(frontendDistFS))
	router.GET("/js/*filepath", serveStaticFile(frontendDistFS))
	router.GET("/img/*filepath", serveStaticFile(frontendDistFS))
	router.GET("/fonts/*filepath", serveStaticFile(frontendDistFS))
	router.GET("/favicon.ico", serveStaticFile(frontendDistFS))

	// Handle other frontend routes (for client-side routing)
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// Skip API routes - they should be handled by their own handlers
		if len(path) >= 4 && path[:4] == "/api" {
			return
		}

		// For all other routes, serve index.html for client-side routing
		serveIndexHTML(c, frontendDistFS)
	})
}

// serveIndexHTML serves the index.html file
func serveIndexHTML(c *gin.Context, fsys fs.FS) {
	indexHTML, err := fs.ReadFile(fsys, "index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to load index.html")
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, string(indexHTML))
}

// serveStaticFile serves static files with proper MIME types
func serveStaticFile(fsys fs.FS) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the file path from the URL
		filePath := c.Request.URL.Path[1:] // Remove leading slash

		// Try to open the file
		file, err := fsys.Open(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				c.String(http.StatusNotFound, "File not found")
				return
			}
			c.String(http.StatusInternalServerError, "Failed to open file")
			return
		}
		defer file.Close()

		// Get file info
		stat, err := file.Stat()
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to get file info")
			return
		}

		// If it's a directory, return 404
		if stat.IsDir() {
			c.String(http.StatusNotFound, "Not a file")
			return
		}

		// Set the content type based on file extension
		ext := filepath.Ext(filePath)
		if contentType := mime.TypeByExtension(ext); contentType != "" {
			c.Header("Content-Type", contentType)
		}

		// Read the file content
		content, err := fs.ReadFile(fsys, filePath)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to read file")
			return
		}

		// Serve the file content
		c.Data(http.StatusOK, c.GetHeader("Content-Type"), content)
	}
}
