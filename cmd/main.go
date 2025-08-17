package main

import (
	"log"
	"net/http"
	"os"
	"wavy-health-api/internal/db"
	"wavy-health-api/internal/handlers"
	"wavy-health-api/internal/services"
)

// withCORS middleware for enabling Cross-Origin Resource Sharing
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Directly get environment variables from the system
	dbURL := os.Getenv("DB_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	initialStockPrice := os.Getenv("INITIAL_STOCK_PRICE")

	if dbURL == "" || jwtSecret == "" || initialStockPrice == "" {
		log.Fatal("Required environment variables not set: DB_URL, JWT_SECRET, INITIAL_STOCK_PRICE")
	}

	// Connect to PostgreSQL database using DB_URL
	db.Connect(dbURL)

	// Convert initialStockPrice to float64 (assuming it's a string in env)
	price, err := strconv.ParseFloat(initialStockPrice, 64)
	if err != nil {
		log.Fatal("Invalid INITIAL_STOCK_PRICE value:", err)
	}

	services.InitStockPrice(price)

	// Register API routes
	http.Handle("/price", withCORS(http.HandlerFunc(handlers.PriceHandler)))
	http.Handle("/buy", withCORS(http.HandlerFunc(handlers.BuyHandler)))
	http.Handle("/sell", withCORS(http.HandlerFunc(handlers.SellHandler)))
	http.Handle("/account", withCORS(http.HandlerFunc(handlers.AccountHandler)))
	http.Handle("/login", withCORS(http.HandlerFunc(handlers.LoginHandler)))
	http.Handle("/logout", withCORS(http.HandlerFunc(handlers.LogoutHandler)))
	http.Handle("/register", withCORS(http.HandlerFunc(handlers.RegisterHandler)))

	// Start server on port from environment or default 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
