package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ErdajtSopjani/LikeMe_API/pkg/handlers"
	"gorm.io/gorm"
)

// RegisterUserInput is the struct for the input of the RegisterUser handler
type RegisterUserInput struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Bio       string `json:"bio"`
}

// RegisterUser is a handler for registering a new user and sending an email confirmation
func RegisterUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the input from the request body
		var input RegisterUserInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			log.Fatal("failed to decode request body:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create a new user in the database
		user := handlers.User{
			Email:     input.Email,
			Username:  input.Username,
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Bio:       input.Bio,
			Token:     GenerateToken(),
		}

		if err := db.Create(&user).Error; err != nil {
			log.Fatal("failed to create user:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send an email confirmation
		// SendConfirmation(user.Email, user.Token)

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
