package helpers

import (
	"log"

	"gorm.io/gorm"
)

// CheckTokenMatch checks if a token and a userId match
func CheckTokenMatch(db *gorm.DB, userId int64, token string) bool {
	var count int64

	// see if the token and userId match in the user_tokens table
	err := db.Table("user_tokens").
		Where("user_id = ? AND token = ?", userId, token).
		Count(&count).Error

	if err != nil {
		log.Printf("Error checking token match: %v", err)
		return false
	}

	// if a record is found the count will go up
	return count > 0
}
