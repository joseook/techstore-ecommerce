package models

type CartItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Cart struct {
	Items    []CartItem `json:"items"`
	Total    float64    `json:"total"`
	Discount float64    `json:"discount"`
}

var UserCart = Cart{
	Items:    []CartItem{},
	Total:    0.0,
	Discount: 0.0,
}

func AddToCart(productID string, quantity int) {
	// Find the product
	var product Product
	for _, p := range Products {
		if p.ID == productID {
			product = p
			break
		}
	}

	// Check if item already exists in cart
	for i, item := range UserCart.Items {
		if item.ProductID == productID {
			UserCart.Items[i].Quantity += quantity
			updateCartTotal()
			return
		}
	}

	// Add new item to cart
	UserCart.Items = append(UserCart.Items, CartItem{
		ProductID: productID,
		Quantity:  quantity,
		Price:     product.Price,
	})

	updateCartTotal()
}

func RemoveFromCart(productID string) {
	for i, item := range UserCart.Items {
		if item.ProductID == productID {
			UserCart.Items = append(UserCart.Items[:i], UserCart.Items[i+1:]...)
			updateCartTotal()
			return
		}
	}
}

func UpdateCartItemQuantity(productID string, quantity int) {
	for i, item := range UserCart.Items {
		if item.ProductID == productID {
			UserCart.Items[i].Quantity = quantity
			updateCartTotal()
			return
		}
	}
}

func updateCartTotal() {
	total := 0.0
	for _, item := range UserCart.Items {
		total += item.Price * float64(item.Quantity)
	}
	UserCart.Total = total - UserCart.Discount
}

func GetCart() Cart {
	return UserCart
}

func ClearCart() {
	UserCart.Items = []CartItem{}
	UserCart.Total = 0.0
	UserCart.Discount = 0.0
} 