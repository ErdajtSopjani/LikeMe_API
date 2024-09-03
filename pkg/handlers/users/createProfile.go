package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ErdajtSopjani/LikeMe_API/pkg/handlers"
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("\n\nBAD REQUEST\n\tBad request: %s\n\tError: %s\n\n", req, err)
			return
		}

		// get userid from the token
		var userToken handlers.UserTokens
		println("running query")
		if err := db.Select("user_id").Where("token = ?", r.Header.Get("Authorization")).First(&userToken).Error; err != nil {
			log.Printf("\n\nERROR\n\tFailed to query database!\n\t%s\n\n", err)
			http.Error(w, "Internal-Server Error...", http.StatusInternalServerError)
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
		if !handlers.CheckUnique(db, "username", userProfile.Username, "user_profile") {
			http.Error(w, "Username already taken", http.StatusBadRequest)
			return
		}

		// check if any field is empty except for bio
		if userProfile.Username == "" || userProfile.FirstName == "" || userProfile.LastName == "" || userProfile.ProfilePicture == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		// create user profile
		if err := db.Create(&userProfile).Error; err != nil { // if profile creation fails
			log.Fatal("failed to create user profile: ", err)
			http.Error(w, "Internal-Server Error...", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User profile created"))
	}
}
