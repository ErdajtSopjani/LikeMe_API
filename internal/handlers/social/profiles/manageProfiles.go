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
type ProfileRequest struct {
	Username       string `json:"username"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture string `json:"profile_picture"`
	Bio            string `json:"bio"`
}

// ManageProfiles is a route for creating and updating user profiles
func ManageProfiles(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ProfileRequest
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
		if r.Method != "PUT" && req.Username == "" || req.FirstName == "" || req.LastName == "" {
			println(r.Method)
			helpers.RespondError(w, "All fields are required", http.StatusBadRequest)
			return
		}

		// check if username is taken by another user
		if db.Where("username = ? AND user_id != ?", userProfile.Username, userProfile.UserId).First(&handlers.UserProfile{}).Error == nil {
			helpers.RespondError(w, "Username already taken", http.StatusBadRequest)
			return
		}

		// check if the user already has a profile
		if !helpers.CheckUnique(db, "user_id", userProfile.UserId, "user_profiles") {
			if r.Method == "PUT" {
				updateProfile(db, userProfile, w)
				return
			}
			helpers.RespondError(w, "User already has a profile, you can update it.", http.StatusBadRequest)
			return
		}

		// create user profile
		if err := db.Create(&userProfile).Error; err != nil { // if profile creation fails
			helpers.RespondError(w, "Internal-Server Error...", http.StatusInternalServerError)
			log.Printf("\n\nERROR\n\tFailed to create user profile: %s\n\n", err)
			return
		}

		helpers.RespondJSON(w, http.StatusCreated, "Profile Created")
	}
}

// updateProfile updates the user profile after verification
func updateProfile(db *gorm.DB, userProfile handlers.UserProfile, w http.ResponseWriter) error {
	if err := db.Model(&handlers.UserProfile{}).Where("user_id = ?", userProfile.UserId).Updates(&userProfile).Error; err != nil {
		helpers.RespondError(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	helpers.RespondJSON(w, http.StatusOK, "Profile Updated")
	return nil
}
