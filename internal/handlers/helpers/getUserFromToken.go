package helpers

import (
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers"
	"gorm.io/gorm"
)

// GetUserFromToken returns the user associated with the given token
func GetUserFromToken(db *gorm.DB, token string) (handlers.User, error) {
	var userToken handlers.UserToken
	err := db.Where("token = ?", token).First(&userToken).Error

	var user handlers.User
	err = db.Where("id = ?", userToken.UserId).First(&user).Error
	return user, err
}
