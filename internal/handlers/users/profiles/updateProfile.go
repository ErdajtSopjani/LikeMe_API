package profiles

import (
	"errors"
	"net/http"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/helpers"
	"gorm.io/gorm"
)

type UpdateProfileRequest struct {
	Username       string `json:"username"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture string `json:"profile_picture"`
	Bio            string `json:"bio"`
}

// updateProfile updates the user profile after verification
func updateProfile(db *gorm.DB, userProfile handlers.UserProfile, w http.ResponseWriter) error {
	// check if username is taken by another user
	if db.Where("username = ? AND user_id != ?", userProfile.Username, userProfile.UserId).First(&handlers.UserProfile{}).Error == nil {
		helpers.RespondError(w, "Username already taken", http.StatusBadRequest)
		return errors.New("Username already taken")
	}

	if err := db.Model(&handlers.UserProfile{}).Where("user_id = ?", userProfile.UserId).Updates(&userProfile).Error; err != nil {
		helpers.RespondError(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	helpers.RespondJSON(w, http.StatusOK, "Profile Updated")
	return nil
}
