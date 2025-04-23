package api

import (
	"fmt"
	"net/http"
)

// Handler is the serverless function entry point for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Golang E-commerce! Path: %s", r.URL.Path)
} 