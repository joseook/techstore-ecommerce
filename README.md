# Golang E-commerce

A modern, lightweight e-commerce platform built with Go and Gin framework.

![Golang E-commerce](https://img.shields.io/badge/Go-E-commerce-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![version](https://img.shields.io/badge/version-1.0.0-blue)

## 📋 Overview

Golang E-commerce is a robust and scalable e-commerce solution built with Go's Gin framework. The application provides essential e-commerce functionality with a clean, responsive UI.

## ✨ Features

- **Product Catalog:** Browse products with category filtering
- **Product Details:** Detailed product information with reviews and related products
- **Shopping Cart:** Add, update, and remove items from your cart
- **Checkout Process:** Simple and secure checkout flow
- **Responsive Design:** Works on desktop and mobile devices
- **REST API:** API endpoints for products, categories, and cart management

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/joseok/golang-ecommerce.git
   cd golang-ecommerce
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the application:
   ```bash
   go run main.go
   ```

4. Open your browser and navigate to `http://localhost:8080`

## 🔧 Tech Stack

- **Backend:** Go, Gin Framework
- **Frontend:** HTML, CSS, JavaScript
- **UI Framework:** Bootstrap

## 📁 Project Structure

```
golang-ecommerce/
├── main.go              # Main application entry point
├── templates/           # HTML templates
│   ├── home.html        # Homepage template
│   ├── products.html    # Products listing page
│   ├── product-detail.html # Product detail page
│   ├── cart.html        # Shopping cart page
│   ├── checkout.html    # Checkout page
│   └── support.html     # Support page
├── static/              # Static assets (JS, CSS, images)
│   ├── css/             # CSS files
│   ├── js/              # JavaScript files
│   └── images/          # Image assets
└── models/              # Data models
```

## 🔄 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/products` | Get all products |
| GET | `/api/products?popular=true` | Get popular products |
| GET | `/api/products?new=true` | Get new products |
| GET | `/api/categories` | Get all categories |
| GET | `/api/brands` | Get all brands |
| POST | `/cart/add` | Add item to cart |
| POST | `/cart/update` | Update cart item |
| POST | `/cart/remove` | Remove item from cart |

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request 