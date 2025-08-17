package main

import (
	"log"
	"net/http"
	"os"
	"wavy-health-api/internal/db"
	"wavy-health-api/internal/handlers"
	"wavy-health-api/internal/services"

	"github.com/joho/godotenv"
)

// withCORS is a simple middleware function that enables
// Cross-Origin Resource Sharing (CORS) for the API endpoints.
// This allows requests from different origins (e.g., your frontend).
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
	// Load environment variables from a .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to the PostgreSQL database
	db.Connect()

	// Initialize stock price with a starting value, e.g., $50
	// This will start the stock price simulation in a separate goroutine
	services.InitStockPrice(50.0)

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
		port = "8080" // Default to 8080 if not specified in .env
	}

	// Start the HTTP server and log its status
	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
