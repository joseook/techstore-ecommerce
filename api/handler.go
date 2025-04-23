package main

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handler is the serverless function entry point for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// Setup Gin in release mode
	gin.SetMode(gin.ReleaseMode)

	// Create a new gin context from the request and writer
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = r

	// Handle static files
	if strings.HasPrefix(r.URL.Path, "/static/") {
		file := strings.TrimPrefix(r.URL.Path, "/static/")
		http.ServeFile(w, r, filepath.Join("static", file))
		return
	}

	// For API demonstration, we'll return JSON data for now
	if r.URL.Path == "/" || r.URL.Path == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"message":  "Golang E-commerce API",
			"products": products,
		})
		return
	}

	// Handle product listing
	if r.URL.Path == "/products" {
		ctx.JSON(http.StatusOK, gin.H{
			"products": products,
		})
		return
	}

	// Handle product details by ID
	if strings.HasPrefix(r.URL.Path, "/product/") {
		idStr := strings.TrimPrefix(r.URL.Path, "/product/")
		var product Product
		found := false

		for _, p := range products {
			if p.ID == idStr {
				product = p
				found = true
				break
			}
		}

		if found {
			ctx.JSON(http.StatusOK, product)
		} else {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		}
		return
	}

	// Handle 404 for all other routes
	ctx.JSON(http.StatusNotFound, gin.H{
		"error": "Not found",
		"path":  r.URL.Path,
	})
}

// Run starts a local development server
func Run() {
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}

