package api

import (
	"net/http"
	"strings"
	"fmt"
	"html/template"

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
	products     []Product
	cart         Cart
	templates    *template.Template
	productCache = make(map[string][]byte)
	// Ensure the main application is initialized only once
	ginApp *gin.Engine
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

// Handler is the entry point for Vercel serverless function
func Handler(w http.ResponseWriter, r *http.Request) {
	// Initialize the Gin application if not already done
	if ginApp == nil {
		gin.SetMode(gin.ReleaseMode)
		ginApp = setupRouter()
	}

	// Extract path from URL
	path := r.URL.Path
	
	// Add support for SPA routing by redirecting to the root path for direct HTML page requests
	if !strings.HasPrefix(path, "/api") && !strings.HasPrefix(path, "/static") && !strings.Contains(path, ".") {
		r.URL.Path = "/"
	}

	// Serve the request
	ginApp.ServeHTTP(w, r)
}

// setupRouter configures and returns the Gin router with all routes
func setupRouter() *gin.Engine {
	// Define modo release para Gin
	gin.SetMode(gin.ReleaseMode)

	// Cria um novo roteador Gin
	r := gin.Default()

	// Configure custom template functions
	r.SetFuncMap(template.FuncMap{
		"iterate": func(count int) []int {
			items := make([]int, count)
			for i := 0; i < count; i++ {
				items[i] = i
			}
			return items
		},
		"multiply": func(a, b float64) float64 {
			return a * b
		},
	})

	// Carrega cada template individualmente
	r.LoadHTMLFiles(
		"templates/home.html",
		"templates/products.html",
		"templates/product-detail.html",
		"templates/cart.html",
		"templates/support.html",
		"templates/checkout.html",
	)

	// Serve arquivos estáticos
	r.Static("/static", "./static")

	// Rotas
	r.GET("/", func(c *gin.Context) {
		// Filter popular and new products for the home page
		var popularProducts []Product
		var newProducts []Product

		for _, product := range products {
			if product.IsPopular {
				popularProducts = append(popularProducts, product)
			}
			if product.IsNew {
				newProducts = append(newProducts, product)
			}
		}
		fmt.Println("Rota / foi chamada! Renderizando home.html...")
		c.HTML(http.StatusOK, "home.html", gin.H{
			"title":           "Home",
			"products":        products,
			"popularProducts": popularProducts,
			"newProducts":     newProducts,
		})
	})

	r.GET("/products", func(c *gin.Context) {
		// Get query parameters
		category := c.Query("category")

		var filteredProducts []Product

		if category != "" {
			// Filter by category
			for _, p := range products {
				if p.Category == category {
					filteredProducts = append(filteredProducts, p)
				}
			}
		} else {
			// Return all products if no specific filter
			filteredProducts = products
		}

		c.HTML(http.StatusOK, "products.html", gin.H{
			"title":    "Our Products",
			"products": filteredProducts,
		})
	})

	r.GET("/product/:id", func(c *gin.Context) {
		productID := c.Param("id")

		var product Product
		found := false

		for _, p := range products {
			if p.ID == productID {
				product = p
				found = true
				break
			}
		}

		if !found {
			c.HTML(http.StatusNotFound, "product-detail.html", gin.H{
				"title":     "Product Not Found",
				"error":     "Product not found",
				"productID": productID,
			})
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

		// Find related products from the same category
		var relatedProducts []Product
		for _, p := range products {
			if p.Category == product.Category && p.ID != product.ID {
				relatedProducts = append(relatedProducts, p)
				// Limit to 3 related products
				if len(relatedProducts) >= 3 {
					break
				}
			}
		}

		c.HTML(http.StatusOK, "product-detail.html", gin.H{
			"title":           product.Name,
			"Product":         product,
			"reviews":         reviews,
			"relatedProducts": relatedProducts,
		})
	})

	// Página de checkout
	r.GET("/checkout", func(c *gin.Context) {
		// Calcular o subtotal e o total do carrinho
		subtotal := 0.0
		for _, item := range cart.Items {
			subtotal += item.Product.Price * float64(item.Quantity)
		}
		
		// Aplicar taxas e frete (simulado)
		tax := subtotal * 0.1  // 10% de imposto
		shipping := 0.0
		if subtotal < 100 {
			shipping = 15.0  // $15 de frete para compras abaixo de $100
		}
		
		total := subtotal + tax + shipping
		
		c.HTML(http.StatusOK, "checkout.html", gin.H{
			"title":     "Checkout",
			"cart":      cart,
			"subtotal":  subtotal,
			"tax":       tax,
			"shipping":  shipping,
			"total":     total,
		})
	})

	// Health check endpoint
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API endpoints for products
	r.GET("/api/products", func(c *gin.Context) {
		// Check for query parameters
		popular := c.Query("popular")
		new := c.Query("new")

		var filteredProducts []Product

		if popular == "true" {
			// Filter popular products
			for _, product := range products {
				if product.IsPopular {
					filteredProducts = append(filteredProducts, product)
				}
			}
		} else if new == "true" {
			// Filter new products
			for _, product := range products {
				if product.IsNew {
					filteredProducts = append(filteredProducts, product)
				}
			}
		} else {
			// Return all products if no specific filter
			filteredProducts = products
		}

		c.JSON(http.StatusOK, filteredProducts)
	})

	// Categories API endpoint
	r.GET("/api/categories", func(c *gin.Context) {
		// Extract unique categories
		categoryMap := make(map[string]bool)
		for _, product := range products {
			categoryMap[product.Category] = true
		}

		// Convert map to slice
		var categories []string
		for category := range categoryMap {
			categories = append(categories, category)
		}

		c.JSON(http.StatusOK, categories)
	})

	// Brands API endpoint
	r.GET("/api/brands", func(c *gin.Context) {
		// For now, we don't have brands in our model, so return dummy data
		brands := []string{"Apple", "Samsung", "Sony", "Microsoft", "Google"}
		c.JSON(http.StatusOK, brands)
	})

	r.GET("/cart", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"cart": cart,
		})
	})

	r.POST("/cart/add", func(c *gin.Context) {
		var request struct {
			Product  struct{ ID string } `json:"product"`
			Quantity int                 `json:"quantity"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Find product by ID
		var product Product
		for _, p := range products {
			if p.ID == request.Product.ID {
				product = p
				break
			}
		}

		if product.ID == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		// Adiciona ao carrinho
		cart.Items = append(cart.Items, CartItem{
			Product:  product,
			Quantity: request.Quantity,
		})

		c.JSON(http.StatusOK, cart)
	})

	r.POST("/cart/update", func(c *gin.Context) {
		var request struct {
			Product  struct{ ID string } `json:"product"`
			Quantity int                 `json:"quantity"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Atualiza quantidade no carrinho
		for i, item := range cart.Items {
			if item.Product.ID == request.Product.ID {
				cart.Items[i].Quantity += request.Quantity
				if cart.Items[i].Quantity <= 0 {
					cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
				}
				break
			}
		}

		c.JSON(http.StatusOK, cart)
	})

	r.POST("/cart/remove", func(c *gin.Context) {
		var request struct {
			Product struct{ ID string } `json:"product"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Remove do carrinho
		for i, item := range cart.Items {
			if item.Product.ID == request.Product.ID {
				cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
				break
			}
		}

		c.JSON(http.StatusOK, cart)
	})

	return r
} 