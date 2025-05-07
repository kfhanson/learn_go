package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// User represents a customer in our store
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"-"` // Never return password in JSON
}

// UserResponse is what we return to clients
type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Token string `json:"token,omitempty"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// In-memory user database for demo purposes
var users = []User{
	{
		ID:       "u1",
		Email:    "john@example.com",
		Name:     "John Doe",
		Password: "password123", // In production, use hashed passwords
	},
}

// Handler processes user-related requests
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
	
	// Handle login endpoint
	if len(pathParts) > 2 && pathParts[2] == "login" {
		handleLogin(w, r)
		return
	}
	
	// Handle register endpoint
	if len(pathParts) > 2 && pathParts[2] == "register" {
		handleRegister(w, r)
		return
	}
	
	// Handle user profile endpoint
	if len(pathParts) > 2 && pathParts[2] != "" {
		userID := pathParts[2]
		handleUserProfile(w, r, userID)
		return
	}
	
	// Method not allowed
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// handleLogin processes login requests
func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Find user by email
	var user *User
	for i := range users {
		if users[i].Email == loginReq.Email {
			user = &users[i]
			break
		}
	}
	
	// User not found or password incorrect
	if user == nil || user.Password != loginReq.Password {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
		return
	}
	
	// Generate a simple token (in production, use JWT)
	token := "token-" + user.ID + "-" + string(time.Now().Unix())
	
	// Return user info with token
	response := UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Token: token,
	}
	
	json.NewEncoder(w).Encode(response)
}

// handleRegister processes registration requests
func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Check if email already exists
	for _, user := range users {
		if user.Email == newUser.Email {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": "Email already in use"})
			return
		}
	}
	
	// Generate a simple ID (in production, use UUID)
	newUser.ID = "u" + string(len(users)+1)
	
	// Add to users
	users = append(users, newUser)
	
	// Generate a simple token (in production, use JWT)
	token := "token-" + newUser.ID + "-" + string(time.Now().Unix())
	
	// Return user info with token
	response := UserResponse{
		ID:    newUser.ID,
		Email: newUser.Email,
		Name:  newUser.Name,
		Token: token,
	}
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// handleUserProfile processes requests for a specific user
func handleUserProfile(w http.ResponseWriter, r *http.Request, id string) {
	// In a real app, verify authentication token here
	
	// Find user by ID
	var user *User
	for i := range users {
		if users[i].ID == id {
			user = &users[i]
			break
		}
	}
	
	// User not found
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
		return
	}
	
	// GET - Return user details
	if r.Method == "GET" {
		response := UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// Method not allowed
	w.WriteHeader(http.StatusMethodNotAllowed)
}
