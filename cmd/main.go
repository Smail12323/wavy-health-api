package main

import (
	"log"
	"net/http"
	"os"
	"wavy-health-api/internal/db"
	"wavy-health-api/internal/handlers"
	"wavy-health-api/internal/services"
)

// withCORS is a simple middleware function that enables
// Cross-Origin Resource Sharing (CORS) for the API endpoints.
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
	// Connect to the PostgreSQL database
	db.Connect()

	// Initialize stock price
	services.InitStockPrice(50.0)

	// Register API routes
	http.Handle("/price", withCORS(http.HandlerFunc(handlers.PriceHandler)))
	http.Handle("/buy", withCORS(http.HandlerFunc(handlers.BuyHandler)))
	http.Handle("/sell", withCORS(http.HandlerFunc(handlers.SellHandler)))
	http.Handle("/account", withCORS(http.HandlerFunc(handlers.AccountHandler)))
	http.Handle("/login", withCORS(http.HandlerFunc(handlers.LoginHandler)))
	http.Handle("/logout", withCORS(http.HandlerFunc(handlers.LogoutHandler)))
	http.Handle("/register", withCORS(http.HandlerFunc(handlers.RegisterHandler)))

	// Get port from environment (Render sets PORT automatically)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
