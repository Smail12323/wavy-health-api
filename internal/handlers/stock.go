package handlers

import (
	"encoding/json"
	"net/http"
	"wavy-health-api/internal/db"
	"wavy-health-api/internal/services"
)

// BuySellRequest represents the expected JSON payload for buying or selling stocks
type BuySellRequest struct {
	UserID int `json:"user_id"` // ID of the user performing the action
	Amount int `json:"amount"`  // Number of stocks to buy or sell
}

// PriceHandler returns the current stock price
func PriceHandler(w http.ResponseWriter, r *http.Request) {
	price := services.GetPrice()
	json.NewEncoder(w).Encode(map[string]float64{"price": price})
}

// BuyHandler handles stock purchases by users
func BuyHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode request JSON into BuySellRequest struct
	var req BuySellRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Calculate total cost of requested stocks
	price := services.GetPrice()
	cost := price * float64(req.Amount)

	// Retrieve user's current balance
	var balance float64
	err := db.DB.QueryRow("SELECT balance FROM users WHERE id=$1", req.UserID).Scan(&balance)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if user has enough balance
	if balance < cost {
		http.Error(w, "Insufficient balance", http.StatusBadRequest)
		return
	}

	// Update user's balance and stock holdings
	_, err = db.DB.Exec(
		"UPDATE users SET balance=balance-$1, stocks_owned=stocks_owned+$2 WHERE id=$3",
		cost, req.Amount, req.UserID,
	)
	if err != nil {
		http.Error(w, "Failed to update account", http.StatusInternalServerError)
		return
	}

	// Return updated account info
	var stocks int
	db.DB.QueryRow("SELECT balance, stocks_owned FROM users WHERE id=$1", req.UserID).Scan(&balance, &stocks)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"balance":      balance,
		"stocks_owned": stocks,
	})
}

// SellHandler handles stock sales by users
func SellHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode request JSON into BuySellRequest struct
	var req BuySellRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	price := services.GetPrice()

	// Retrieve user's current stock holdings
	var stocks int
	err := db.DB.QueryRow("SELECT stocks_owned FROM users WHERE id=$1", req.UserID).Scan(&stocks)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if user has enough stocks to sell
	if stocks < req.Amount {
		http.Error(w, "Not enough stocks", http.StatusBadRequest)
		return
	}

	// Update user's balance and stock holdings
	_, err = db.DB.Exec(
		"UPDATE users SET balance=balance+$1, stocks_owned=stocks_owned-$2 WHERE id=$3",
		price*float64(req.Amount), req.Amount, req.UserID,
	)
	if err != nil {
		http.Error(w, "Failed to update account", http.StatusInternalServerError)
		return
	}

	// Return updated account info
	var balance float64
	db.DB.QueryRow("SELECT balance, stocks_owned FROM users WHERE id=$1", req.UserID).Scan(&balance, &stocks)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"balance":      balance,
		"stocks_owned": stocks,
	})
}
