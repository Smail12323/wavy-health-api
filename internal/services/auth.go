package services

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret is the secret key used to sign JWTs, loaded from environment variables
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// GenerateJWT generates a JWT token for a given user ID
// The token includes the user ID and expires in 24 hours
func GenerateJWT(userID int) (string, error) {
	// Create a new JWT token with HS256 signing method and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,                                // Store user ID in token claims
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Set token expiration to 24 hours
	})

	// Sign the token with the secret key and return it
	return token.SignedString(jwtSecret)
}
