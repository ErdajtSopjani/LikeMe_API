package profiles

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
	Username       string `json:"username"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture string `json:"profile_picture"`
	Bio            string `json:"bio"`
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
		var userToken handlers.UserToken
		if err := db.Select("user_id").Where("token = ?", r.Header.Get("Authorization")).First(&userToken).Error; err != nil {
			log.Printf("\n\nERROR\n\tFailed to query database!\n\t%s\n\n", err)
			helpers.RespondError(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		userProfile := handlers.UserProfile{
			Username:       req.Username,
			FirstName:      req.FirstName,
			LastName:       req.LastName,
			ProfilePicture: req.ProfilePicture,
			Bio:            req.Bio,
			UserId:         userToken.UserId,
		}

		// check if any field is empty except for bio
		if req.Username == "" || req.FirstName == "" || req.LastName == "" {
			helpers.RespondError(w, "All fields are required", http.StatusBadRequest)
			return
		}

		// check if the user already has a profile
		if !helpers.CheckUnique(db, "user_id", userProfile.UserId, "user_profiles") {
			helpers.RespondError(w, "User already has a profile", http.StatusBadRequest)
			return
		}

		// check if the username is taken
		if !helpers.CheckUnique(db, "username", userProfile.Username, "user_profiles") {
			helpers.RespondError(w, "Username already taken", http.StatusBadRequest)
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
