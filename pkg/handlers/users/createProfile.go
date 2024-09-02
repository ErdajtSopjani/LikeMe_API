package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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

// UserProfile represents the structure of user_profile table in the database
type UserProfile struct {
	Username       string `json:"username"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture string `json:"profile_picture"`
	Bio            string `json:"bio"`
	UserId         int64  `json:"user_id"`
}

// UserTokens represents the structure of user_tokens table in the database
type UserTokens struct {
	Id        int64
	ExpiresAt time.Time
	CreatedAt time.Time
	UserId    int64
}

// CreateProfile creates the user profile after verification
func CreateProfile(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateProfileRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Fatal("failed to decode request body: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// get userid from the token
		var userToken UserTokens
		if err := db.Select("user_id").Where("token = ?", r.Header.Get("Authorization")).First(&userToken).Error; err != nil {
			log.Fatal("failed to query database: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("user with token found: ", r.Header.Get("Authorization"), " ", userToken.UserId)

		userProfile := UserProfile{
			Username:       req.Username,
			FirstName:      req.FirstName,
			LastName:       req.LastName,
			ProfilePicture: req.ProfilePicture,
			Bio:            req.Bio,
			UserId:         userToken.UserId,
		}

		// check if the username is taken
		if !handlers.CheckUnique(db, "username", userProfile.Username) {
			http.Error(w, "Username already taken", http.StatusBadRequest)
			return
		}

		// create user profile
		if err := db.Create(&userProfile).Error; err != nil { // if profile creation fails
			log.Fatal("failed to create user profile: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User profile created"))
	}
}
