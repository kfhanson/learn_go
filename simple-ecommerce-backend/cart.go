package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

// CartItem represents an item in a user's cart
type CartItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

// Cart represents a user's shopping cart
type Cart struct {
	UserID string     `json:"userId"`
	Items  []CartItem `json:"items"`
}

// CartRequest represents a request to add/update cart items
type CartRequest struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

// In-memory cart database for demo purposes
var carts = []Cart{
	{
		UserID: "u1",
		Items: []CartItem{
			{ProductID: "p1", Quantity: 2},
		},
	},
}

// Handler processes cart-related requests
func Handler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	
	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// In a real app, get userID from authentication token
	// For demo, we'll extract it from the URL
	path := r.URL.Path
	pathParts := strings.Split(path, "/")
	
	if len(pathParts) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "User ID required"})
		return
	}
	
	userID := pathParts[2]
	
	// Handle different methods
	switch r.Method {
	case "GET":
		getCart(w, userID)
	case "POST":
		addToCart(w, r, userID)
	case "PUT":
		updateCart(w, r, userID)
	case "DELETE":
		removeFromCart(w, r, userID)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// getCart returns the user's cart
func getCart(w http.ResponseWriter, userID string) {
	// Find cart by userID
	var cart *Cart
	for i := range carts {
		if carts[i].UserID == userID {
			cart = &carts[i]
			break
		}
	}
	
	// If cart doesn't exist, create an empty one
	if cart == nil {
		cart = &Cart{
			UserID: userID,
			Items:  []CartItem{},
		}
		carts = append(carts, *cart)
	}
	
	json.NewEncoder(w).Encode(cart)
}

// addToCart adds an item to the cart
func addToCart(w http.ResponseWriter, r *http.Request, userID string) {
	var req CartRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate quantity
	if req.Quantity <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Quantity must be positive"})
		return
	}
	
	// Find cart by userID
	var cart *Cart
	var cartIndex int
	for i := range carts {
		if carts[i].UserID == userID {
			cart = &carts[i]
			cartIndex = i
			break
		}
	}
	
	// If cart doesn't exist, create a new one
	if cart == nil {
		cart = &Cart{
			UserID: userID,
			Items:  []CartItem{},
		}
		carts = append(carts, *cart)
		cartIndex = len(carts) - 1
	}
	
	// Check if product already in cart
	var found bool
	for i := range cart.Items {
		if cart.Items[i].ProductID == req.ProductID {
			// Update quantity
			cart.Items[i].Quantity += req.Quantity
			found = true
			break
		}
	}
	
	// If product not in cart, add it
	if !found {
		cart.Items = append(cart.Items, CartItem{
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		})
	}
	
	// Update cart in storage
	carts[cartIndex] = *cart
	
	json.NewEncoder(w).Encode(cart)
}

// updateCart updates the quantity of an item in the cart
func updateCart(w http.ResponseWriter, r *http.Request, userID string) {
	var req CartRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Find cart by userID
	var cart *Cart
	var cartIndex int
	for i := range carts {
		if carts[i].UserID == userID {
			cart = &carts[i]
			cartIndex = i
			break
		}
	}
	
	// Cart not found
	if cart == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cart not found"})
		return
	}
	
	// Find product in cart
	var found bool
	for i := range cart.Items {
		if cart.Items[i].ProductID == req.ProductID {
			// If quantity is 0 or negative, remove item
			if req.Quantity <= 0 {
				cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			} else {
				// Update quantity
				cart.Items[i].Quantity = req.Quantity
			}
			found = true
			break
		}
	}
	
	// Product not in cart
	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Product not in cart"})
		return
	}
	
	// Update cart in storage
	carts[cartIndex] = *cart
	
	json.NewEncoder(w).Encode(cart)
}

// removeFromCart removes an item from the cart
func removeFromCart(w http.ResponseWriter, r *http.Request, userID string) {
	// Extract product ID from query parameters
	productID := r.URL.Query().Get("productId")
	if productID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Product ID required"})
		return
	}
	
	// Find cart by userID
	var cart *Cart
	var cartIndex int
	for i := range carts {
		if carts[i].UserID == userID {
			cart = &carts[i]
			cartIndex = i
			break
		}
	}
	
	// Cart not found
	if cart == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cart not found"})
		return
	}
	
	// Find product in cart
	var found bool
	for i := range cart.Items {
		if cart.Items[i].ProductID == productID {
			// Remove item
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			found = true
			break
		}
	}
	
	// Product not in cart
	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Product not in cart"})
		return
	}
	
	// Update cart in storage
	carts[cartIndex] = *cart
	
	json.NewEncoder(w).Encode(cart)
}
