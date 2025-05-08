package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// OrderItem represents an item in an order
type OrderItem struct {
	ProductID string  `json:"productId"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}

// Order represents a customer order
type Order struct {
	ID           string      `json:"id"`
	UserID       string      `json:"userId"`
	Items        []OrderItem `json:"items"`
	TotalAmount  float64     `json:"totalAmount"`
	Status       string      `json:"status"`
	CreatedAt    time.Time   `json:"createdAt"`
	ShippingAddr Address     `json:"shippingAddress"`
}

// Address represents a shipping address
type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zipCode"`
	Country string `json:"country"`
}

// OrderRequest represents a request to create an order
type OrderRequest struct {
	ShippingAddr Address `json:"shippingAddress"`
}

// In-memory order database for demo purposes
var orders = []Order{
	{
		ID:          "o1",
		UserID:      "u1",
		TotalAmount: 259.98,
		Status:      "processing",
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		Items: []OrderItem{
			{
				ProductID: "p1",
				Name:      "Mechanical Keyboard",
				Price:     129.99,
				Quantity:  2,
			},
		},
		ShippingAddr: Address{
			Street:  "123 Main St",
			City:    "Anytown",
			State:   "CA",
			ZipCode: "12345",
			Country: "USA",
		},
	},
}

// Handler processes order-related requests
func Handler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	
	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Extract path parts
	path := r.URL.Path
	pathParts := strings.Split(path, "/")
	
	// In a real app, get userID from authentication token
	// For demo, we'll extract it from the URL if present
	if len(pathParts) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "User ID required"})
		return
	}
	
	userID := pathParts[2]
	
	// Handle specific order
	if len(pathParts) > 3 && pathParts[3] != "" {
		orderID := pathParts[3]
		getOrder(w, userID, orderID)
		return
	}
	
	// Handle different methods
	switch r.Method {
	case "GET":
		getOrders(w, userID)
	case "POST":
		createOrder(w, r, userID)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// getOrders returns all orders for a user
func getOrders(w http.ResponseWriter, userID string) {
	// Find orders for user
	var userOrders []Order
	for _, order := range orders {
		if order.UserID == userID {
			userOrders = append(userOrders, order)
		}
	}
	
	json.NewEncoder(w).Encode(userOrders)
}

// getOrder returns a specific order
func getOrder(w http.ResponseWriter, userID, orderID string) {
	// Find order by ID
	var order *Order
	for i := range orders {
		if orders[i].ID == orderID {
			order = &orders[i]
			break
		}
	}
	
	// Order not found
	if order == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Order not found"})
		return
	}
	
	// Verify order belongs to user
	if order.UserID != userID {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "Access denied"})
		return
	}
	
	json.NewEncoder(w).Encode(order)
}

// createOrder creates a new order from the user's cart
func createOrder(w http.ResponseWriter, r *http.Request, userID string) {
	var req OrderRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Find user's cart
	var cart *Cart
	for i := range carts {
		if carts[i].UserID == userID {
			cart = &carts[i]
			break
		}
	}
	
	// Cart not found or empty
	if cart == nil || len(cart.Items) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cart is empty"})
		return
	}
	
	// Create order items and calculate total
	var orderItems []OrderItem
	var totalAmount float64
	
	for _, item := range cart.Items {
		// Find product details
		var product *Product
		for i := range products {
			if products[i].ID == item.ProductID {
				product = &products[i]
				break
			}
		}
		
		if product == nil {
			continue // Skip if product not found
		}
		
		// Check stock
		if product.Stock < item.Quantity {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Not enough stock for " + product.Name,
			})
			return
		}
		
		// Add to order items
		orderItems = append(orderItems, OrderItem{
			ProductID: product.ID,
			Name:      product.Name,
			Price:     product.Price,
			Quantity:  item.Quantity,
		})
		
		// Update total
		totalAmount += product.Price * float64(item.Quantity)
		
		// Update stock
		product.Stock -= item.Quantity
	}
	
	// Create new order
	newOrder := Order{
		ID:           "o" + string(len(orders)+1),
		UserID:       userID,
		Items:        orderItems,
		TotalAmount:  totalAmount,
		Status:       "pending",
		CreatedAt:    time.Now(),
		ShippingAddr: req.ShippingAddr,
	}
	
	// Add to orders
	orders = append(orders, newOrder)
	
	// Clear cart
	for i := range carts {
		if carts[i].UserID == userID {
			carts[i].Items = []CartItem{}
			break
		}
	}
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newOrder)
}
