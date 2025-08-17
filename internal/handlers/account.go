package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wavy-health-api/internal/db"
)

// AccountHandler handles requests to get a user's account information,
// including current balance and number of stocks owned.
// Expects a query parameter "user_id" in the URL.
func AccountHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user_id from the URL query
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		// Return an error if the user ID is invalid
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Variables to hold user data
	var balance float64
	var stocks int

	// Query the database for the user's balance and stocks_owned
	err = db.DB.QueryRow("SELECT balance, stocks_owned FROM users WHERE id=$1", userID).Scan(&balance, &stocks)
	if err != nil {
		// Return an error if the user is not found
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Respond with JSON containing balance and stocks_owned
	json.NewEncoder(w).Encode(map[string]interface{}{
		"balance":      balance,
		"stocks_owned": stocks,
	})
}
