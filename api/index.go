package api

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// Product representa um produto na loja
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

// CartItem representa um item no carrinho
type CartItem struct {
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
}

// Cart representa o carrinho de compras
type Cart struct {
	Items []CartItem `json:"items"`
}

// Variáveis globais
var (
	products []Product
	cart     Cart
)

func init() {
	// Inicializa produtos
	products = []Product{
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
		{
			ID:          "3",
			Name:        "Wireless Earbuds",
			Description: "True wireless earbuds with long battery life",
			Price:       129.99,
			Image:       "https://img.freepik.com/free-psd/wireless-earbuds-mockup_53876-58593.jpg",
			ImageURL:    "https://img.freepik.com/free-psd/wireless-earbuds-mockup_53876-58593.jpg",
			Category:    "Electronics",
			IsNew:       false,
			IsPopular:   true,
		},
		{
			ID:          "4",
			Name:        "Bluetooth Speaker",
			Description: "Portable waterproof speaker with 20-hour battery",
			Price:       89.99,
			Image:       "https://img.freepik.com/free-psd/bluetooth-speaker-mockup_53876-58594.jpg",
			ImageURL:    "https://img.freepik.com/free-psd/bluetooth-speaker-mockup_53876-58594.jpg",
			Category:    "Electronics",
			IsNew:       false,
			IsPopular:   false,
		},
	}

	// Inicializa carrinho
	cart = Cart{
		Items: make([]CartItem, 0),
	}
}

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