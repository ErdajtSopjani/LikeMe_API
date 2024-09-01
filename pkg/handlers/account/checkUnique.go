package handlers

import "gorm.io/gorm"

// CheckUnique checks if any value is unique in a table
func CheckUnique(db *gorm.DB, column string, value any) bool {
	type User struct {
		Id int64
	}
	var user User

	/// query database to check if the value is unique
	// if no record is found, database returns an error
	if err := db.Select("id").Where(column+" = ?", value).First(&user).Error; err == nil {
		return false
	}

	// if database doesnt return an error that means the record was found
	return true
}
