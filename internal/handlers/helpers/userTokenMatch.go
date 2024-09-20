package helpers

import (
	"log"

	"gorm.io/gorm"
)

// CheckTokenMatch checks if a token and a userId match
func CheckTokenMatch(db *gorm.DB, userId int64, token string) bool {
	var count int64

	// see if the token and userId match in the user_tokens table
	if err := db.Table("user_tokens").Where("user_id = ? AND token = ?", userId, token).Count(&count).Error; err != nil {
		log.Panicf("failed to check if token matches user: %v", err)
		return false
	}

	return count == 0
}
