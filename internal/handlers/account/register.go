package account

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/email"
	"gorm.io/gorm"
)

// RegisterUser is a handler for registering a new user and sending an email confirmation
func RegisterUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the req from the request body
		var req handlers.User
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("\n\nBAD REQUEST\n\tBad request: %v\n\tError: %s\n\n", req, err)
			return
		}

		// Create a new user in the database
		user := &handlers.User{
			Email:       req.Email,
			CountryCode: req.CountryCode,
		}

		// check if the email exists in the database
		if !handlers.CheckUnique(db, "email", user.Email, "users") { // if email exists
			http.Error(w, "Email already taken", http.StatusBadRequest)
			return
		}

		// create the user
		if err := db.Create(&user).Error; err != nil {
			http.Error(w, "Internal Server Error...", http.StatusInternalServerError)
			log.Printf("\n\nERROR\n\tFailed to create user:%v\n\tError: %s\n\n", user, err)
			return
		}

		// get the user id from the database
		// since it's needed for both sending the email and for creating the token
		var userId int64
		if err := db.Table("users").Select("id").Where("email = ?", user.Email).Scan(&userId).Error; err != nil {
			http.Error(w, "Internal Server Error.", http.StatusInternalServerError)
			log.Printf("\n\nERROR\n\tFailed to get user id: %v\n\tError: %s", req, err)
			return
		}

		// create and save a token for the user
		userTokens := &handlers.VerificationTokens{
			UserId: userId,
			Token:  handlers.GenerateToken(),
		}

		if err := db.Create(&userTokens).Error; err != nil || userTokens.Token == "" {
			// if the token is not saved or created
			http.Error(w, "Internal Server Error...", http.StatusInternalServerError)
			log.Printf("\n\nERROR\n\tFailed to create or save token for user: %v\n\tError: %s\n\n", user, err)
			return
		}

		// send the confirmation email
		if err := email.SendConfirmation(db, user.Email, userId); err != nil {
			http.Error(w, "Internal Server Error.", http.StatusInternalServerError)
			log.Printf("\n\nERROR\n\tFailed to send confirmation email: %v\n\tError: %s\n\n", req, err)
			return
		}

		// Respond with the new user
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User created"))
	}
}
