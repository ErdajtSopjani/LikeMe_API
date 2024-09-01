package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ErdajtSopjani/LikeMe_API/pkg/handlers"
	"gorm.io/gorm"
)

// RegisterRequest is the struct for the req of the RegisterUser handler
type RegisterRequest struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Bio       string `json:"bio"`
}

// RegisterUser is a handler for registering a new user and sending an email confirmation
func RegisterUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the req from the request body
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Fatal("failed to decode req:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create a new user in the database
		user := handlers.User{
			Email:     req.Email,
			Username:  req.Username,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Bio:       req.Bio,
			Token:     GenerateToken(),
		}

		if !CheckUnique(db, "username", user.Username) { // if username exists
			http.Error(w, "Username already taken", http.StatusBadRequest)
			return
		}
		if !CheckUnique(db, "email", user.Email) { // if email exists
			http.Error(w, "Email already taken", http.StatusBadRequest)
			return
		}

		if err := db.Create(&user).Error; err != nil {
			log.Fatal("failed to create user:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send an email confirmation
		// email.SendConfirmation(user.Email)

		// Respond with the new user
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			log.Fatal("failed to encode response:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
