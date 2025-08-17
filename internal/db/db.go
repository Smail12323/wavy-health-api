package db

import (
	"database/sql"
	"log"
	"os"

	// PostgreSQL driver
	_ "github.com/lib/pq"
)

// DB is the global database connection object
var DB *sql.DB

// Connect initializes the connection to the PostgreSQL database.
// It reads the database URL from the environment variable DB_URL,
// opens the connection, and verifies it by pinging the database.
func Connect() {
	// Get the database URL from environment variables
	dbURL := os.Getenv("DB_URL")
	log.Println("Connecting to DB with URL:", dbURL)

	var err error

	// Open a connection to the database
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cannot open DB:", err)
	}

	// Verify the connection is alive
	if err = DB.Ping(); err != nil {
		log.Fatal("Cannot ping DB:", err)
	}

	log.Println("Connected to PostgreSQL successfully")
}
