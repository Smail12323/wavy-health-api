package handlers

import (
	"encoding/json"
	"net/http"
	"wavy-health-api/internal/db"
	"wavy-health-api/internal/services"

	"golang.org/x/crypto/bcrypt"
)

// LoginRequest represents the expected JSON payload for a login request.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the JSON response sent after a successful login.
type LoginResponse struct {
	Token string `json:"token"` // JWT token to authenticate future requests
	ID    int    `json:"id"`    // User ID
}

// LoginHandler handles user login requests.
// It expects a POST request with JSON containing username and password.
// If credentials are valid, it returns a JWT token and user ID in JSON format.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON body into a LoginRequest struct
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Fetch the user's ID and password hash from the database
	var id int
	var hash string
	err = db.DB.QueryRow("SELECT id, password_hash FROM users WHERE username=$1", req.Username).Scan(&id, &hash)
	if err != nil {
		// If user is not found or query fails, return unauthorized
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare the provided password with the stored hash
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate a JWT token for the authenticated user
	token, err := services.GenerateJWT(id)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Prepare and send the JSON response
	resp := LoginResponse{
		Token: token,
		ID:    id,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
