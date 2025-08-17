package handlers

import (
	"net/http"
)

// LogoutHandler handles user logout requests.
// Since authentication is done via JWT, logout is mainly client-side.
// The client should delete the stored token to "log out" the user.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Return OK status indicating that logout action is complete
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out"))
}
