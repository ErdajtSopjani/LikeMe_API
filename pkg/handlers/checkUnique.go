package handlers

import (
	"gorm.io/gorm"
	"log"
)

// CheckUnique checks if any value is unique in a table
func CheckUnique(db *gorm.DB, column string, value any) bool {
	type User struct {
		Id int64
	}
	var user User

	// check for the value in the database
	if err := db.Select("id").Where(column+" = ?", value).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound { // if the record is not found an error gets returned
			return true
		}
		// if the error is not a not found error
		log.Fatalf("Error checking uniqueness for %s: %v\n", column, err)
		return false
	}

	// if no error is returned it means the value was found
	return false
}
