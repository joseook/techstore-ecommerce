package main

import (
	"fmt"
	"github.com/joseook/techstore-ecommerce/models"
)

func main() {
	fmt.Println("Testing imports...")
	fmt.Printf("Models loaded: %d products available\n", len(models.Products))
} 