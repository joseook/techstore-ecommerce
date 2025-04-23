package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"github.com/joseok/techstore-ecommerce/models"
)

type CartItem struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}

type Cart struct {
	Items []CartItem `json:"items"`
	Total float64    `json:"total"`
}

var (
	cart     = Cart{
		Items: make([]CartItem, 0),
		Total: 0,
	}
	cartLock sync.Mutex
)

func GetCart(w http.ResponseWriter, r *http.Request) {
	cartLock.Lock()
	defer cartLock.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

func AddToCart(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Quantity <= 0 {
		http.Error(w, "Quantity must be greater than 0", http.StatusBadRequest)
		return
	}

	// Get product details from models
	product, err := getProductByID(request.ProductID)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	cartLock.Lock()
	defer cartLock.Unlock()

	// Update or add item to cart
	found := false
	for i, item := range cart.Items {
		if item.ProductID == request.ProductID {
			cart.Items[i].Quantity += request.Quantity
			found = true
			break
		}
	}

	if !found {
		cart.Items = append(cart.Items, CartItem{
			ProductID: product.ID,
			Name:      product.Name,
			Price:     product.Price,
			Quantity:  request.Quantity,
		})
	}

	// Update total
	cart.Total = 0
	for _, item := range cart.Items {
		cart.Total += item.Price * float64(item.Quantity)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cart)
}

func UpdateCart(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Quantity <= 0 {
		http.Error(w, "Quantity must be greater than 0", http.StatusBadRequest)
		return
	}

	cartLock.Lock()
	defer cartLock.Unlock()

	found := false
	for i, item := range cart.Items {
		if item.ProductID == request.ProductID {
			cart.Items[i].Quantity = request.Quantity
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Item not found in cart", http.StatusNotFound)
		return
	}

	// Update total
	cart.Total = 0
	for _, item := range cart.Items {
		cart.Total += item.Price * float64(item.Quantity)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cart)
}

func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Path[len("/api/cart/"):]

	cartLock.Lock()
	defer cartLock.Unlock()

	found := false
	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Item not found in cart", http.StatusNotFound)
		return
	}

	// Update total
	cart.Total = 0
	for _, item := range cart.Items {
		cart.Total += item.Price * float64(item.Quantity)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cart)
}

func getProductByID(id string) (*models.Product, error) {
	for _, product := range models.Products {
		if product.ID == id {
			return &product, nil
		}
	}
	return nil, fmt.Errorf("product not found")
} 