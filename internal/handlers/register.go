package handlers

import (
	"encoding/json"
	"net/http"
	"wavy-health-api/internal/db"

	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest represents the expected JSON payload for user registration
type RegisterRequest struct {
	Username string `json:"username"` // Desired username
	Password string `json:"password"` // Plain-text password
}

// RegisterHandler handles new user registrations.
// It hashes the password and stores the user in the database.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode JSON request body into RegisterRequest struct
	var req RegisterRequest
	json.NewDecoder(r.Body).Decode(&req)

	// Hash the password before storing it
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	// Insert new user into database
	_, err := db.DB.Exec(
		"INSERT INTO users (username, password_hash) VALUES ($1, $2)",
		req.Username,
		string(hash),
	)
	if err != nil {
		// Return error if username already exists or any DB error occurs
		http.Error(w, "User already exists or DB error", http.StatusBadRequest)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered"})
}
