package api

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// Product represents a product in the store
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Image       string  `json:"image"`
	ImageURL    string  `json:"image_url"`
	Category    string  `json:"category"`
	IsNew       bool    `json:"isNew"`
	IsPopular   bool    `json:"isPopular"`
}

// Products is a list of sample products
var Products = []Product{
	{
		ID:          "1",
		Name:        "Premium Headphones",
		Description: "High-quality wireless headphones with noise cancellation",
		Price:       199.99,
		Image:       "https://img.freepik.com/free-psd/headphones-mockup_53876-58591.jpg",
		ImageURL:    "https://img.freepik.com/free-psd/headphones-mockup_53876-58591.jpg",
		Category:    "Electronics",
		IsNew:       true,
		IsPopular:   true,
	},
	{
		ID:          "2",
		Name:        "Smart Watch",
		Description: "Feature-rich smartwatch with health monitoring",
		Price:       249.99,
		Image:       "https://img.freepik.com/free-psd/smartwatch-mockup_53876-58592.jpg",
		ImageURL:    "https://img.freepik.com/free-psd/smartwatch-mockup_53876-58592.jpg",
		Category:    "Electronics",
		IsNew:       true,
		IsPopular:   false,
	},
}

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
			"products": Products,
		})
		return
	}

	// Handle product listing
	if r.URL.Path == "/products" {
		ctx.JSON(http.StatusOK, gin.H{
			"products": Products,
		})
		return
	}

	// Handle product details by ID
	if strings.HasPrefix(r.URL.Path, "/product/") {
		idStr := strings.TrimPrefix(r.URL.Path, "/product/")
		var product Product
		found := false

		for _, p := range Products {
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