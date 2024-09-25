package account

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/email/verify"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/helpers"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Email       string `json:"email"`
	CountryCode string `json:"country_code"`
}

// RegisterUser is a handler for registering a new user and sending an email confirmation
func RegisterUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the req from the request body
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helpers.RespondError(w, "Invalid Request-Format", http.StatusBadRequest)
			log.Printf("\n\nBAD REQUEST\n\tBad request: %v\n\tError: %s\n\n", req, err)
			return
		}

		// validate email format using regex
		emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		if !regexp.MustCompile(emailRegex).MatchString(req.Email) {
			helpers.RespondError(w, "Invalid Email", http.StatusBadRequest)
			log.Printf("Invalid email format: %v\n", req.Email)
			return
		}

		// check if the country_code is missing
		if req.CountryCode == "" {
			helpers.RespondError(w, "Country Code is required", http.StatusBadRequest)
			log.Printf("Missing Country Code")
			return
		}

		// create a new user structure
		user := &handlers.User{
			Email:       req.Email,
			CountryCode: req.CountryCode,
			CreatedAt:   &handlers.Now,
		}

		// check if the email exists in the database
		if !helpers.CheckUnique(db, "email", user.Email, "users") { // if email exists
			helpers.RespondError(w, "Email already taken", http.StatusBadRequest)
			return
		}

		// save the user to the database
		if err := db.Create(&user).Error; err != nil {
			helpers.RespondError(w, err, http.StatusInternalServerError)
			log.Printf("\n\nERROR\n\tFailed to create user:%v\n\tError: %s\n\n", user, err)
			return
		}

		// get the user id from the database
		// since it's needed for both sending the email and for creating the token
		var userId int64
		if err := db.Table("users").Select("id").Where("email = ?", user.Email).Scan(&userId).Error; err != nil {
			helpers.RespondError(w, "Internal Server Error.", http.StatusInternalServerError)
			log.Printf("\n\nERROR\n\tFailed to get user id: %v\n\tError: %s", req, err)
			return
		}

		// send the confirmation email
		// NOTE: this will also create and save a token available for 10 minutes that the user can use to verify their email
		if err := verify.SendConfirmation(db, user.Email, userId); err != nil {
			helpers.RespondError(w, "Internal Server Error.", http.StatusInternalServerError)
			log.Printf("\n\nERROR\n\tFailed to send confirmation email: %v\n\tError: %s\n\n", req, err)
			return
		}

		helpers.RespondJSON(w, http.StatusCreated, "User created")
	}
}
