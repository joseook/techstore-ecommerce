package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	// Ensure the main application is initialized only once
	ginApp *gin.Engine
)

// Handle is the serverless function handler for Vercel
func Handle(w http.ResponseWriter, r *http.Request) {
	// Initialize the Gin application if not already done
	if ginApp == nil {
		gin.SetMode(gin.ReleaseMode)
		ginApp = setupRouter()
	}

	// Extract path from URL
	path := r.URL.Path
	
	// Add support for SPA routing by redirecting to the root path for direct HTML page requests
	if !strings.HasPrefix(path, "/api") && !strings.HasPrefix(path, "/static") && !strings.Contains(path, ".") {
		r.URL.Path = "/"
	}

	// Serve the request
	ginApp.ServeHTTP(w, r)
}

// This is a simplified version of the setupRouter function from main.go
// You should maintain the same routes as in main.go
func setupRouter() *gin.Engine {
	// Import the real setupRouter from main package
	// For Vercel, we need to include all routes in this file
	// This is a placeholder - copy the router setup from main.go
	
	r := gin.Default()
	
	// Basic route to test the serverless function
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	
	// Add other routes from main.go here
	// ...
	
	return r
} 