package users

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/helpers"
	"gorm.io/gorm"
)

// CreateProfileRequest is the structure of the request body for the createProfile endpoint
type CreateProfileRequest struct {
	Username       string
	FirstName      string
	LastName       string
	ProfilePicture string
	Bio            string
}

// CreateProfile creates the user profile after verification
func CreateProfile(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateProfileRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helpers.RespondError(w, err.Error(), http.StatusBadRequest)
			log.Printf("\n\nBAD REQUEST\n\tBad request: %s\n\tError: %s\n\n", req, err)
			return
		}

		// get userid from the token
		var userToken handlers.VerificationToken
		println("running query")
		if err := db.Select("user_id").Where("token = ?", r.Header.Get("Authorization")).First(&userToken).Error; err != nil {
			log.Printf("\n\nERROR\n\tFailed to query database!\n\t%s\n\n", err)
			helpers.RespondError(w, "Internal-Server Error...", http.StatusInternalServerError)
			return
		}

		log.Println("user with token found: ", r.Header.Get("Authorization"), " ", userToken.UserId)

		userProfile := handlers.UserProfile{
			Username:       req.Username,
			FirstName:      req.FirstName,
			LastName:       req.LastName,
			ProfilePicture: req.ProfilePicture,
			Bio:            req.Bio,
			UserId:         userToken.UserId,
		}

		// check if the username is taken
		if !helpers.CheckUnique(db, "username", userProfile.Username, "user_profile") {
			helpers.RespondError(w, "Username already taken", http.StatusBadRequest)
			return
		}

		// check if any field is empty except for bio
		if userProfile.Username == "" || userProfile.FirstName == "" || userProfile.LastName == "" || userProfile.ProfilePicture == "" {
			helpers.RespondError(w, "All fields are required", http.StatusBadRequest)
			return
		}

		// create user profile
		if err := db.Create(&userProfile).Error; err != nil { // if profile creation fails
			helpers.RespondError(w, "Internal-Server Error...", http.StatusInternalServerError)
			log.Printf("\n\nERROR\n\tFailed to create user profile: %s\n\n", err)
			return
		}

		helpers.RespondJSON(w, http.StatusCreated, "User profile created")
	}
}
