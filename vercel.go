package handler

import (
	"net/http"
	
	"github.com/joseook/techstore-ecommerce/api"
)

// Handler is a wrapper function that calls the api.Handler
func Handler(w http.ResponseWriter, r *http.Request) {
	api.Handler(w, r)
}

