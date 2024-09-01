package handlers

import (
	"gorm.io/gorm"
	"log"
)

func UserExists(db *gorm.DB, userId int64) bool {
	type User struct {
		id int64
	}
	var user User

	// query database to check if the user with (userId) exists
	err := db.Select("id").Where("id = ?", userId).First(&user).Error
	if err != nil {
		log.Println("User with id: |", userId, "| doesn't exist!")
		return false
	}

	return true
}
