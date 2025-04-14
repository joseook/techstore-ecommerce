package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
	"sort"
	"github.com/joseok/golang-ecommerce/models"
	"github.com/gorilla/mux"
)

// Templates map to store parsed templates
var templates map[string]*template.Template

// InitTemplates initializes the templates map
func InitTemplates(t map[string]*template.Template) {
	templates = t
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Get query parameters
	category := r.URL.Query().Get("category")
	brand := r.URL.Query().Get("brand")
	search := r.URL.Query().Get("search")
	popular := r.URL.Query().Get("popular")
	new := r.URL.Query().Get("new")
	sort := r.URL.Query().Get("sort")

	// Filter products
	filteredProducts := models.Products

	if category != "" {
		filteredProducts = filterByCategory(filteredProducts, category)
	}

	if brand != "" {
		filteredProducts = filterByBrand(filteredProducts, brand)
	}

	if search != "" {
		filteredProducts = filterBySearch(filteredProducts, search)
	}

	if popular == "true" {
		filteredProducts = filterPopularProducts(filteredProducts)
	}

	if new == "true" {
		filteredProducts = filterNewProducts(filteredProducts)
	}
	
	// Sort the filtered products
	if sort != "" {
		filteredProducts = sortProducts(filteredProducts, sort)
	}

	json.NewEncoder(w).Encode(filteredProducts)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	for _, product := range models.Products {
		if product.ID == id {
			json.NewEncoder(w).Encode(product)
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	categories := make(map[string]bool)
	for _, product := range models.Products {
		categories[product.Category] = true
	}

	var result []string
	for category := range categories {
		result = append(result, category)
	}

	json.NewEncoder(w).Encode(result)
}

func GetBrands(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	brands := make(map[string]bool)
	for _, product := range models.Products {
		brands[product.Brand] = true
	}

	var result []string
	for brand := range brands {
		result = append(result, brand)
	}

	json.NewEncoder(w).Encode(result)
}

func filterByCategory(products []models.Product, category string) []models.Product {
	var filtered []models.Product
	for _, product := range products {
		if product.Category == category {
			filtered = append(filtered, product)
		}
	}
	return filtered
}

func filterByBrand(products []models.Product, brand string) []models.Product {
	var filtered []models.Product
	for _, product := range products {
		if product.Brand == brand {
			filtered = append(filtered, product)
		}
	}
	return filtered
}

func filterBySearch(products []models.Product, search string) []models.Product {
	var filtered []models.Product
	for _, product := range products {
		if strings.Contains(strings.ToLower(product.Name), strings.ToLower(search)) ||
			strings.Contains(strings.ToLower(product.Description), strings.ToLower(search)) {
			filtered = append(filtered, product)
		}
	}
	return filtered
}

func filterPopularProducts(products []models.Product) []models.Product {
	var filtered []models.Product
	for _, product := range products {
		if product.IsPopular {
			filtered = append(filtered, product)
		}
	}
	return filtered
}

func filterNewProducts(products []models.Product) []models.Product {
	var filtered []models.Product
	for _, product := range products {
		if product.IsNew {
			filtered = append(filtered, product)
		}
	}
	return filtered
}

// GetProductsByCategory returns products filtered by category
func GetProductsByCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	category := r.URL.Query().Get("category")
	var filteredProducts []models.Product
	for _, product := range models.Products {
		if strings.ToLower(product.Category) == strings.ToLower(category) {
			filteredProducts = append(filteredProducts, product)
		}
	}
	json.NewEncoder(w).Encode(filteredProducts)
}

// GetProductsByBrand returns products filtered by brand
func GetProductsByBrand(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	brand := r.URL.Query().Get("brand")
	var filteredProducts []models.Product
	for _, product := range models.Products {
		if strings.ToLower(product.Brand) == strings.ToLower(brand) {
			filteredProducts = append(filteredProducts, product)
		}
	}
	json.NewEncoder(w).Encode(filteredProducts)
}

// SearchProducts returns products matching the search query
func SearchProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := strings.ToLower(r.URL.Query().Get("q"))
	var filteredProducts []models.Product
	for _, product := range models.Products {
		if strings.Contains(strings.ToLower(product.Name), query) ||
			strings.Contains(strings.ToLower(product.Description), query) {
			filteredProducts = append(filteredProducts, product)
		}
	}
	json.NewEncoder(w).Encode(filteredProducts)
}

// GetProductDetail returns the product detail page
func GetProductDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	var product models.Product
	for _, p := range models.Products {
		if p.ID == productID {
			product = p
			break
		}
	}

	if product.ID == "" {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Sample reviews (in a real app, these would come from a database)
	reviews := []struct {
		User    string
		Rating  float64
		Comment string
		Date    string
	}{
		{
			User:    "John Doe",
			Rating:  5,
			Comment: "Great product! Very satisfied with my purchase.",
			Date:    "2024-03-15",
		},
		{
			User:    "Jane Smith",
			Rating:  4,
			Comment: "Good quality, but a bit expensive.",
			Date:    "2024-03-10",
		},
	}

	data := struct {
		Title   string
		Product models.Product
		Reviews []struct {
			User    string
			Rating  float64
			Comment string
			Date    string
		}
	}{
		Title:   "TechStore - " + product.Name,
		Product: product,
		Reviews: reviews,
	}

	if err := templates["product-detail.html"].Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// sortProducts sorts the products based on the sort option
func sortProducts(products []models.Product, sortOption string) []models.Product {
	sortedProducts := make([]models.Product, len(products))
	copy(sortedProducts, products)
	
	switch sortOption {
	case "popular":
		// Sort by popularity (IsPopular first, then by rating)
		sort.Slice(sortedProducts, func(i, j int) bool {
			if sortedProducts[i].IsPopular != sortedProducts[j].IsPopular {
				return sortedProducts[i].IsPopular
			}
			return sortedProducts[i].Rating > sortedProducts[j].Rating
		})
	case "newest":
		// Sort by newness (IsNew first)
		sort.Slice(sortedProducts, func(i, j int) bool {
			return sortedProducts[i].IsNew
		})
	case "price-low":
		// Sort by price (low to high)
		sort.Slice(sortedProducts, func(i, j int) bool {
			return sortedProducts[i].Price < sortedProducts[j].Price
		})
	case "price-high":
		// Sort by price (high to low)
		sort.Slice(sortedProducts, func(i, j int) bool {
			return sortedProducts[i].Price > sortedProducts[j].Price
		})
	}
	
	return sortedProducts
} 