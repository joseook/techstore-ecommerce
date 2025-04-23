package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler is the main entry point for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// Create a new Gin router
	router := setupRouter()
	
	// Convert the http.Request to a Gin context
	ginContext := &gin.Context{
		Request: r,
		Writer:  w,
	}
	
	// Handle the request
	router.HandleContext(ginContext)
}

// SetupRouter sets up the Gin router
func setupRouter() *gin.Engine {
	router := gin.Default()
	
	// Handle static files
	router.Static("/static", "static")
	
	// For API demonstration, we'll return JSON data for now
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":  "Golang E-commerce API",
			"products": products,
		})
	})
	
	// Handle product listing
	router.GET("/products", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"products": products,
		})
	})
	
	// Handle product details by ID
	router.GET("/product/:id", func(c *gin.Context) {
		id := c.Param("id")
		var product Product
		found := false
		
		for _, p := range products {
			if p.ID == id {
				product = p
				found = true
				break
			}
		}
		
		if found {
			c.JSON(http.StatusOK, product)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		}
	})
	
	return router
} 