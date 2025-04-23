package api

import (
	"encoding/json"
	"net/http"
)

// Product representa um produto na loja
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Image       string  `json:"image"`
	Category    string  `json:"category"`
}

// Simple products database
var products = []Product{
	{
		ID:          "1",
		Name:        "Premium Headphones",
		Description: "High-quality wireless headphones with noise cancellation",
		Price:       199.99,
		Image:       "https://img.freepik.com/free-psd/headphones-mockup_53876-58591.jpg",
		Category:    "Electronics",
	},
	{
		ID:          "2",
		Name:        "Smart Watch",
		Description: "Feature-rich smartwatch with health monitoring",
		Price:       249.99,
		Image:       "https://img.freepik.com/free-psd/smartwatch-mockup_53876-58592.jpg",
		Category:    "Electronics",
	},
}

// Handler is the Vercel serverless function
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Return API status and products
	response := map[string]interface{}{
		"status": "ok",
		"message": "Golang E-commerce API is running",
		"products": products,
	}
	
	json.NewEncoder(w).Encode(response)
} 