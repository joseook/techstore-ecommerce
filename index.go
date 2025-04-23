package main

import (
	"fmt"
	"net/http"

	"github.com/joseook/techstore-ecommerce/pkg/api"
)

// Handler is the main entry point for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	api.Handler(w, r)
}

// main function for local development
func main() {
	fmt.Println("Starting server on :8080...")
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
} 