package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Product represents an item in our store
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"imageUrl"`
	Stock       int     `json:"stock"`
}

// In-memory product database for demo purposes
// In a real application, you would use a database
var products = []Product{
	{
		ID:          "p1",
		Name:        "Mechanical Keyboard",
		Description: "Premium mechanical keyboard with RGB lighting",
		Price:       129.99,
		ImageURL:    "https://example.com/keyboard.jpg",
		Stock:       50,
	},
	{
		ID:          "p2",
		Name:        "Wireless Mouse",
		Description: "Ergonomic wireless mouse with long battery life",
		Price:       49.99,
		ImageURL:    "https://example.com/mouse.jpg",
		Stock:       100,
	},
	{
		ID:          "p3",
		Name:        "Monitor Stand",
		Description: "Adjustable monitor stand for better ergonomics",
		Price:       79.99,
		ImageURL:    "https://example.com/stand.jpg",
		Stock:       30,
	},
}

// Handler processes product-related requests
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

	// Extract product ID from path if present
	path := r.URL.Path
	pathParts := strings.Split(path, "/")
	
	// Handle different endpoints
	if len(pathParts) > 2 && pathParts[2] != "" {
		productID := pathParts[2]
		handleSingleProduct(w, r, productID)
		return
	}
	
	// List all products
	if r.Method == "GET" {
		json.NewEncoder(w).Encode(products)
		return
	}
	
	// Create a new product
	if r.Method == "POST" {
		var newProduct Product
		err := json.NewDecoder(r.Body).Decode(&newProduct)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		
		// Generate a simple ID (in production, use UUID)
		newProduct.ID = "p" + string(len(products)+1)
		
		// Add to products
		products = append(products, newProduct)
		
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newProduct)
		return
	}
	
	// Method not allowed
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// handleSingleProduct handles requests for a specific product
func handleSingleProduct(w http.ResponseWriter, r *http.Request, id string) {
	// Find product by ID
	var product *Product
	for i := range products {
		if products[i].ID == id {
			product = &products[i]
			break
		}
	}
	
	// Product not found
	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Product not found"})
		return
	}
	
	// GET - Return product details
	if r.Method == "GET" {
		json.NewEncoder(w).Encode(product)
		return
	}
	
	// PUT - Update product
	if r.Method == "PUT" {
		var updatedProduct Product
		err := json.NewDecoder(r.Body).Decode(&updatedProduct)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		
		// Preserve ID
		updatedProduct.ID = product.ID
		
		// Update product
		for i := range products {
			if products[i].ID == id {
				products[i] = updatedProduct
				break
			}
		}
		
		json.NewEncoder(w).Encode(updatedProduct)
		return
	}
	
	// DELETE - Remove product
	if r.Method == "DELETE" {
		var newProducts []Product
		for i := range products {
			if products[i].ID != id {
				newProducts = append(newProducts, products[i])
			}
		}
		products = newProducts
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted"})
		return
	}
	
	// Method not allowed
	w.WriteHeader(http.StatusMethodNotAllowed)
}
