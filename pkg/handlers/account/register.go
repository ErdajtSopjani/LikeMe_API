package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ErdajtSopjani/LikeMe_API/pkg/handlers"
	"gorm.io/gorm"
)

// RegisterUser is a handler for registering a new user and sending an email confirmation
func RegisterUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the req from the request body
		var req handlers.User
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Fatal("failed to decode req:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create a new user in the database
		user := handlers.User{
			Email:       req.Email,
			CountryCode: req.CountryCode,
		}

		if !handlers.CheckUnique(db, "email", user.Email, handlers.UserProfile{}) { // if email exists
			http.Error(w, "Email already taken", http.StatusBadRequest)
			return
		}

		if err := db.Create(&user).Error; err != nil {
			log.Fatal("failed to create user:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO: Send an email confirmation
		// email.SendConfirmation(user.Email)

		// Respond with the new user
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User created"))
	}
}
