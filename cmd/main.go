package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"wavy-health-api/internal/db"
	"wavy-health-api/internal/handlers"
	"wavy-health-api/internal/services"
)

// withCORS is a simple middleware function that enables
// Cross-Origin Resource Sharing (CORS) for the API endpoints.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	// 🔹 NO .env file, use Render environment variables
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	// Connect to the PostgreSQL database
	db.Connect(dbURL)

	// Get initial stock price from env (or fallback to 50.0)
	initialStockPrice := 50.0
	if val := os.Getenv("INITIAL_STOCK_PRICE"); val != "" {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			initialStockPrice = f
		}
	}

	// Initialize stock price
	services.InitStockPrice(initialStockPrice)

	// Register API routes and wrap them with CORS middleware
	http.Handle("/price", withCORS(http.HandlerFunc(handlers.PriceHandler)))
	http.Handle("/buy", withCORS(http.HandlerFunc(handlers.BuyHandler)))
	http.Handle("/sell", withCORS(http.HandlerFunc(handlers.SellHandler)))
	http.Handle("/account", withCORS(http.HandlerFunc(handlers.AccountHandler)))
	http.Handle("/login", withCORS(http.HandlerFunc(handlers.LoginHandler)))
	http.Handle("/logout", withCORS(http.HandlerFunc(handlers.LogoutHandler)))
	http.Handle("/register", withCORS(http.HandlerFunc(handlers.RegisterHandler)))

	// Determine the port to run the server on
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the HTTP server and log its status
	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
