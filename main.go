package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

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

// Função para verificar ambiente Vercel
func isVercelEnvironment() bool {
	return os.Getenv("VERCEL") != "" || os.Getenv("VERCEL_ENV") != ""
}

// Templates embutidos para ambiente Vercel
var (
	// Templates básicos para ambiente Vercel
	homeTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>TechStore - Home</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; }
        header { background-color: #333; color: white; padding: 10px; text-align: center; }
        .products { display: flex; flex-wrap: wrap; gap: 20px; margin-top: 20px; }
        .product { border: 1px solid #ddd; padding: 15px; width: 250px; }
        img { max-width: 100%; height: auto; }
        .price { font-weight: bold; color: #e44d26; }
        .btn { background-color: #4CAF50; color: white; padding: 10px 15px; text-decoration: none; display: inline-block; margin-top: 10px; }
    </style>
</head>
<body>
    <header>
        <h1>Welcome to TechStore</h1>
        <nav>
            <a href="/">Home</a> |
            <a href="/products">Products</a>
        </nav>
    </header>
    
    <main>
        <h2>Popular Products</h2>
        <div class="products">
            {{range .popularProducts}}
            <div class="product">
                <img src="{{.ImageURL}}" alt="{{.Name}}">
                <h3>{{.Name}}</h3>
                <p class="price">${{.Price}}</p>
                <p>{{.Description}}</p>
                <a href="/product/{{.ID}}" class="btn">View Details</a>
            </div>
            {{end}}
        </div>
        
        <h2>New Arrivals</h2>
        <div class="products">
            {{range .newProducts}}
            <div class="product">
                <img src="{{.ImageURL}}" alt="{{.Name}}">
                <h3>{{.Name}}</h3>
                <p class="price">${{.Price}}</p>
                <p>{{.Description}}</p>
                <a href="/product/{{.ID}}" class="btn">View Details</a>
            </div>
            {{end}}
        </div>
    </main>
    
    <footer style="margin-top: 50px; text-align: center; color: #777;">
        &copy; 2024 TechStore. All rights reserved.
    </footer>
</body>
</html>
`

	productDetailTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.title}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; }
        header { background-color: #333; color: white; padding: 10px; text-align: center; }
        .product-container { display: flex; gap: 30px; margin-top: 30px; }
        .product-image { flex: 1; }
        .product-info { flex: 1; }
        img { max-width: 100%; height: auto; }
        .price { font-size: 24px; font-weight: bold; color: #e44d26; margin: 15px 0; }
        .btn { background-color: #4CAF50; color: white; padding: 10px 15px; text-decoration: none; display: inline-block; border: none; cursor: pointer; font-size: 16px; }
        .reviews { margin-top: 40px; }
        .review { border-bottom: 1px solid #eee; padding: 15px 0; }
        .rating { color: gold; }
    </style>
</head>
<body>
    <header>
        <h1>TechStore</h1>
        <nav>
            <a href="/">Home</a> |
            <a href="/products">Products</a>
        </nav>
    </header>
    
    <main>
        {{if .error}}
            <h2>Error: {{.error}}</h2>
            <p>Product with ID {{.productID}} was not found.</p>
            <a href="/products" class="btn">Back to Products</a>
        {{else}}
            <div class="product-container">
                <div class="product-image">
                    <img src="{{.Product.ImageURL}}" alt="{{.Product.Name}}">
                </div>
                <div class="product-info">
                    <h1>{{.Product.Name}}</h1>
                    <p class="price">${{.Product.Price}}</p>
                    <p>{{.Product.Description}}</p>
                    <p><strong>Category:</strong> {{.Product.Category}}</p>
                    
                    <button class="btn">Add to Cart</button>
                </div>
            </div>
            
            <div class="reviews">
                <h2>Customer Reviews</h2>
                {{range .reviews}}
                <div class="review">
                    <p class="rating">
                        {{range (iterate 5)}}
                            {{if lt . (int $.Rating)}}★{{else}}☆{{end}}
                        {{end}}
                    </p>
                    <p>"{{.Comment}}"</p>
                    <p><strong>{{.User}}</strong> - {{.Date}}</p>
                </div>
                {{end}}
            </div>
            
            <div style="margin-top: 40px;">
                <h2>Related Products</h2>
                <div class="products" style="display: flex; gap: 20px;">
                    {{range .relatedProducts}}
                    <div style="border: 1px solid #ddd; padding: 15px; width: 250px;">
                        <img src="{{.ImageURL}}" alt="{{.Name}}" style="max-width: 100%;">
                        <h3>{{.Name}}</h3>
                        <p class="price">${{.Price}}</p>
                        <a href="/product/{{.ID}}" class="btn">View Details</a>
                    </div>
                    {{end}}
                </div>
            </div>
        {{end}}
    </main>
    
    <footer style="margin-top: 50px; text-align: center; color: #777;">
        &copy; 2024 TechStore. All rights reserved.
    </footer>
</body>
</html>
`

	productsTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.title}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; }
        header { background-color: #333; color: white; padding: 10px; text-align: center; }
        .container { display: flex; gap: 30px; margin-top: 30px; }
        .sidebar { width: 250px; }
        .products { display: flex; flex-wrap: wrap; gap: 20px; flex: 1; }
        .product { border: 1px solid #ddd; padding: 15px; width: 250px; }
        img { max-width: 100%; height: auto; }
        .price { font-weight: bold; color: #e44d26; }
        .btn { background-color: #4CAF50; color: white; padding: 10px 15px; text-decoration: none; display: inline-block; margin-top: 10px; }
        .filter-section { margin-bottom: 20px; border: 1px solid #eee; padding: 15px; }
        .filter-title { font-weight: bold; margin-bottom: 10px; }
        .filter-option { margin: 5px 0; }
    </style>
</head>
<body>
    <header>
        <h1>TechStore - Products</h1>
        <nav>
            <a href="/">Home</a> |
            <a href="/products">Products</a>
        </nav>
    </header>
    
    <main>
        <div class="container">
            <div class="sidebar">
                <div class="filter-section">
                    <div class="filter-title">Categories</div>
                    <div class="filter-option"><a href="/products">All Products</a></div>
                    <div class="filter-option"><a href="/products?category=Electronics">Electronics</a></div>
                </div>
                
                <div class="filter-section">
                    <div class="filter-title">Price Range</div>
                    <div class="filter-option"><a href="/products">All Prices</a></div>
                    <div class="filter-option"><a href="/products">Under $100</a></div>
                    <div class="filter-option"><a href="/products">$100 - $200</a></div>
                    <div class="filter-option"><a href="/products">Over $200</a></div>
                </div>
            </div>
            
            <div>
                <h2>{{.title}}</h2>
                <div class="products">
                    {{range .products}}
                    <div class="product">
                        <img src="{{.ImageURL}}" alt="{{.Name}}">
                        <h3>{{.Name}}</h3>
                        <p class="price">${{.Price}}</p>
                        <p>{{.Description}}</p>
                        <a href="/product/{{.ID}}" class="btn">View Details</a>
                    </div>
                    {{end}}
                </div>
            </div>
        </div>
    </main>
    
    <footer style="margin-top: 50px; text-align: center; color: #777;">
        &copy; 2024 TechStore. All rights reserved.
    </footer>
</body>
</html>
`

	cartTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Shopping Cart</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; }
        header { background-color: #333; color: white; padding: 10px; text-align: center; }
        .cart-container { margin-top: 30px; }
        .cart-item { display: flex; border-bottom: 1px solid #eee; padding: 15px 0; align-items: center; }
        .item-image { width: 100px; margin-right: 20px; }
        .item-details { flex: 1; }
        .item-actions { width: 150px; text-align: right; }
        .quantity-control { display: flex; align-items: center; justify-content: flex-end; margin-bottom: 10px; }
        .btn-quantity { width: 30px; height: 30px; font-size: 16px; }
        .quantity { padding: 0 10px; }
        .summary { margin-top: 30px; padding: 20px; background-color: #f9f9f9; }
        .summary-row { display: flex; justify-content: space-between; margin: 10px 0; }
        .total { font-size: 20px; font-weight: bold; }
        .btn { background-color: #4CAF50; color: white; padding: 10px 15px; text-decoration: none; display: inline-block; border: none; cursor: pointer; font-size: 16px; }
        .btn-checkout { background-color: #e44d26; width: 100%; margin-top: 20px; padding: 15px; font-size: 18px; }
    </style>
</head>
<body>
    <header>
        <h1>Your Shopping Cart</h1>
        <nav>
            <a href="/">Home</a> |
            <a href="/products">Continue Shopping</a>
        </nav>
    </header>
    
    <main>
        <div class="cart-container">
            <h2>Cart Items</h2>
            
            {{if .cart.Items}}
                {{range .cart.Items}}
                <div class="cart-item">
                    <img src="{{.Product.ImageURL}}" alt="{{.Product.Name}}" class="item-image">
                    <div class="item-details">
                        <h3>{{.Product.Name}}</h3>
                        <p class="price">${{.Product.Price}}</p>
                    </div>
                    <div class="item-actions">
                        <div class="quantity-control">
                            <button class="btn-quantity">-</button>
                            <span class="quantity">{{.Quantity}}</span>
                            <button class="btn-quantity">+</button>
                        </div>
                        <button class="btn">Remove</button>
                    </div>
                </div>
                {{end}}
                
                <div class="summary">
                    <div class="summary-row">
                        <span>Subtotal:</span>
                        <span>$100.00</span>
                    </div>
                    <div class="summary-row">
                        <span>Tax (10%):</span>
                        <span>$10.00</span>
                    </div>
                    <div class="summary-row">
                        <span>Shipping:</span>
                        <span>$5.00</span>
                    </div>
                    <div class="summary-row total">
                        <span>Total:</span>
                        <span>$115.00</span>
                    </div>
                    
                    <a href="/checkout" class="btn btn-checkout">Proceed to Checkout</a>
                </div>
            {{else}}
                <p>Your cart is empty. <a href="/products">Start shopping</a></p>
            {{end}}
        </div>
    </main>
    
    <footer style="margin-top: 50px; text-align: center; color: #777;">
        &copy; 2024 TechStore. All rights reserved.
    </footer>
</body>
</html>
`

	checkoutTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Checkout</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; }
        header { background-color: #333; color: white; padding: 10px; text-align: center; }
        .checkout-container { display: flex; gap: 30px; margin-top: 30px; }
        .checkout-form { flex: 2; }
        .order-summary { flex: 1; background-color: #f9f9f9; padding: 20px; height: fit-content; }
        .form-group { margin-bottom: 20px; }
        label { display: block; margin-bottom: 5px; font-weight: bold; }
        input, select { width: 100%; padding: 10px; font-size: 16px; border: 1px solid #ddd; }
        .form-row { display: flex; gap: 15px; }
        .form-row .form-group { flex: 1; }
        .summary-row { display: flex; justify-content: space-between; margin: 10px 0; }
        .total { font-size: 20px; font-weight: bold; margin-top: 15px; padding-top: 15px; border-top: 1px solid #ddd; }
        .btn { background-color: #4CAF50; color: white; padding: 15px; text-decoration: none; display: inline-block; border: none; cursor: pointer; font-size: 18px; width: 100%; margin-top: 20px; }
        .order-item { display: flex; margin-bottom: 10px; }
        .item-details { margin-left: 10px; flex: 1; }
        .item-image { width: 50px; height: 50px; object-fit: cover; }
    </style>
</head>
<body>
    <header>
        <h1>Checkout</h1>
        <nav>
            <a href="/">Home</a> |
            <a href="/products">Products</a> |
            <a href="/cart">Cart</a>
        </nav>
    </header>
    
    <main>
        <div class="checkout-container">
            <div class="checkout-form">
                <h2>Billing Details</h2>
                <form>
                    <div class="form-row">
                        <div class="form-group">
                            <label for="firstName">First Name</label>
                            <input type="text" id="firstName" required>
                        </div>
                        <div class="form-group">
                            <label for="lastName">Last Name</label>
                            <input type="text" id="lastName" required>
                        </div>
                    </div>
                    
                    <div class="form-group">
                        <label for="email">Email Address</label>
                        <input type="email" id="email" required>
                    </div>
                    
                    <div class="form-group">
                        <label for="address">Street Address</label>
                        <input type="text" id="address" required>
                    </div>
                    
                    <div class="form-row">
                        <div class="form-group">
                            <label for="city">City</label>
                            <input type="text" id="city" required>
                        </div>
                        <div class="form-group">
                            <label for="state">State</label>
                            <input type="text" id="state" required>
                        </div>
                        <div class="form-group">
                            <label for="zip">ZIP Code</label>
                            <input type="text" id="zip" required>
                        </div>
                    </div>
                    
                    <h2>Payment Information</h2>
                    
                    <div class="form-group">
                        <label for="cardNumber">Card Number</label>
                        <input type="text" id="cardNumber" placeholder="1234 5678 9012 3456" required>
                    </div>
                    
                    <div class="form-row">
                        <div class="form-group">
                            <label for="expDate">Expiration Date</label>
                            <input type="text" id="expDate" placeholder="MM/YY" required>
                        </div>
                        <div class="form-group">
                            <label for="cvv">CVV</label>
                            <input type="text" id="cvv" required>
                        </div>
                    </div>
                    
                    <button type="submit" class="btn">Place Order</button>
                </form>
            </div>
            
            <div class="order-summary">
                <h2>Order Summary</h2>
                
                {{range .cart.Items}}
                <div class="order-item">
                    <img src="{{.Product.ImageURL}}" alt="{{.Product.Name}}" class="item-image">
                    <div class="item-details">
                        <div>{{.Product.Name}} x {{.Quantity}}</div>
                        <div>${{multiply .Product.Price (float64 .Quantity)}}</div>
                    </div>
                </div>
                {{end}}
                
                <div class="summary-row">
                    <span>Subtotal:</span>
                    <span>${{.subtotal}}</span>
                </div>
                <div class="summary-row">
                    <span>Tax:</span>
                    <span>${{.tax}}</span>
                </div>
                <div class="summary-row">
                    <span>Shipping:</span>
                    <span>${{.shipping}}</span>
                </div>
                <div class="summary-row total">
                    <span>Total:</span>
                    <span>${{.total}}</span>
                </div>
            </div>
        </div>
    </main>
    
    <footer style="margin-top: 50px; text-align: center; color: #777;">
        &copy; 2024 TechStore. All rights reserved.
    </footer>
</body>
</html>
`

	supportTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Customer Support</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; }
        header { background-color: #333; color: white; padding: 10px; text-align: center; }
        .support-container { max-width: 800px; margin: 30px auto; }
        .faq-item { margin-bottom: 20px; border-bottom: 1px solid #eee; padding-bottom: 20px; }
        .faq-question { font-weight: bold; margin-bottom: 10px; font-size: 18px; }
        .contact-form { margin-top: 40px; background-color: #f9f9f9; padding: 20px; }
        .form-group { margin-bottom: 20px; }
        label { display: block; margin-bottom: 5px; font-weight: bold; }
        input, textarea { width: 100%; padding: 10px; font-size: 16px; border: 1px solid #ddd; }
        textarea { height: 150px; }
        .btn { background-color: #4CAF50; color: white; padding: 10px 15px; text-decoration: none; display: inline-block; border: none; cursor: pointer; font-size: 16px; }
    </style>
</head>
<body>
    <header>
        <h1>Customer Support</h1>
        <nav>
            <a href="/">Home</a> |
            <a href="/products">Products</a>
        </nav>
    </header>
    
    <main>
        <div class="support-container">
            <h2>Frequently Asked Questions</h2>
            
            <div class="faq-item">
                <div class="faq-question">How do I track my order?</div>
                <div class="faq-answer">
                    Once your order ships, you will receive a shipping confirmation email with a tracking number. You can use this number to track your package on the carrier's website.
                </div>
            </div>
            
            <div class="faq-item">
                <div class="faq-question">What is your return policy?</div>
                <div class="faq-answer">
                    We accept returns within 30 days of purchase. Items must be in original condition with all packaging and accessories. Please contact our support team to initiate a return.
                </div>
            </div>
            
            <div class="faq-item">
                <div class="faq-question">How long does shipping take?</div>
                <div class="faq-answer">
                    Standard shipping typically takes 3-5 business days. Express shipping options are available at checkout for faster delivery.
                </div>
            </div>
            
            <div class="faq-item">
                <div class="faq-question">Do you offer international shipping?</div>
                <div class="faq-answer">
                    Yes, we ship to most countries worldwide. International shipping rates and delivery times vary by location.
                </div>
            </div>
            
            <div class="contact-form">
                <h2>Contact Us</h2>
                <p>Can't find the answer you're looking for? Please fill out the form below.</p>
                
                <form>
                    <div class="form-group">
                        <label for="name">Your Name</label>
                        <input type="text" id="name" required>
                    </div>
                    
                    <div class="form-group">
                        <label for="email">Email Address</label>
                        <input type="email" id="email" required>
                    </div>
                    
                    <div class="form-group">
                        <label for="orderNumber">Order Number (if applicable)</label>
                        <input type="text" id="orderNumber">
                    </div>
                    
                    <div class="form-group">
                        <label for="subject">Subject</label>
                        <input type="text" id="subject" required>
                    </div>
                    
                    <div class="form-group">
                        <label for="message">Message</label>
                        <textarea id="message" required></textarea>
                    </div>
                    
                    <button type="submit" class="btn">Submit</button>
                </form>
            </div>
        </div>
    </main>
    
    <footer style="margin-top: 50px; text-align: center; color: #777;">
        &copy; 2024 TechStore. All rights reserved.
    </footer>
</body>
</html>
`
)

// setupRouter configures and returns the Gin router
func setupRouter() *gin.Engine {
	// Define modo release para Gin
	gin.SetMode(gin.ReleaseMode)

	// Cria um novo roteador Gin
	r := gin.Default()

	// Log do ambiente atual
	dir, _ := os.Getwd()
	log.Printf("Current working directory: %s", dir)
	log.Printf("Listing files in current directory:")
	files, _ := os.ReadDir(dir)
	for _, file := range files {
		log.Printf("- %s", file.Name())
	}
	
	// Ver se a pasta templates existe
	log.Printf("Checking templates directory:")
	templatesDir := filepath.Join(dir, "templates")
	templateFiles, err := os.ReadDir(templatesDir)
	if err != nil {
		log.Printf("Error reading templates directory: %v", err)
	} else {
		for _, file := range templateFiles {
			log.Printf("- Template: %s", file.Name())
		}
	}

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

	// Verifica se estamos no ambiente Vercel
	if isVercelEnvironment() {
		log.Printf("Using embedded templates for Vercel environment")
		
		// Parse templates from strings
		tmpl := template.New("templates").Funcs(r.FuncMap)
		template.Must(tmpl.New("home.html").Parse(homeTemplate))
		template.Must(tmpl.New("products.html").Parse(productsTemplate))
		template.Must(tmpl.New("product-detail.html").Parse(productDetailTemplate))
		template.Must(tmpl.New("cart.html").Parse(cartTemplate))
		template.Must(tmpl.New("checkout.html").Parse(checkoutTemplate))
		template.Must(tmpl.New("support.html").Parse(supportTemplate))
		
		r.SetHTMLTemplate(tmpl)
	} else {
		// Carrega templates usando caminhos absolutos para ambiente local
		templatePath := filepath.Join(dir, "templates")
		log.Printf("Loading templates from: %s", templatePath)
		r.LoadHTMLFiles(
			filepath.Join(templatePath, "home.html"),
			filepath.Join(templatePath, "products.html"),
			filepath.Join(templatePath, "product-detail.html"),
			filepath.Join(templatePath, "cart.html"),
			filepath.Join(templatePath, "support.html"),
			filepath.Join(templatePath, "checkout.html"),
		)
	}

	// Serve arquivos estáticos com caminho absoluto
	staticPath := filepath.Join(dir, "static")
	log.Printf("Serving static files from: %s", staticPath)
	r.Static("/static", staticPath)

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

	// Adicione uma rota para servir produtos do modelo models/product.go
	// Esta é uma rota de backup para quando a API principal não funcionar
	r.GET("/products/api", func(c *gin.Context) {
		// Importar dinamicamente os produtos do modelo
		var modelProducts []map[string]interface{}
		
		// Converter produtos do main.go para formato compatível com a API
		for _, p := range products {
			product := map[string]interface{}{
				"id":          p.ID,
				"name":        p.Name,
				"description": p.Description,
				"price":       p.Price,
				"image":       p.Image,
				"image_url":   p.ImageURL,
				"category":    p.Category,
				"isNew":       p.IsNew,
				"isPopular":   p.IsPopular,
				// Adicione campos compatíveis com o modelo Product.go
				"is_new":     p.IsNew,
				"is_popular": p.IsPopular,
				"brand":      "Unknown", // Valor padrão
				"rating":     4.5,       // Valor padrão
			}
			modelProducts = append(modelProducts, product)
		}
		
		c.JSON(http.StatusOK, modelProducts)
	})

	return r
}

// Handler is the main entry point for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// Log para debug
	log.Printf("Vercel Handler called with path: %s", r.URL.Path)
	
	// Verificar se é uma solicitação para um recurso estático que pode ser servido diretamente
	if strings.HasPrefix(r.URL.Path, "/static/") || strings.HasPrefix(r.URL.Path, "/templates/") {
		http.FileServer(http.Dir(".")).ServeHTTP(w, r)
		return
	}
	
	// Verificar ambiente Vercel
	if os.Getenv("VERCEL") != "" {
		log.Printf("Running in Vercel environment")
	}
	
	// Create a new Gin router
	router := setupRouter()
	
	// Create a Gin context using CreateTestContext to properly wrap the ResponseWriter
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = r
	
	// Handle the request
	router.HandleContext(ginContext)
}

func main() {
	// Check if we're running in Vercel
	if os.Getenv("VERCEL") != "" {
		// In Vercel, we don't need to start the server
		// The Handler function will be used instead
		return
	}

	// For local development
	router := setupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}