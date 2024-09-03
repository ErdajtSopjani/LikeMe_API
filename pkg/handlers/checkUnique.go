package handlers

import (
	"log"

	"gorm.io/gorm"
)

// CheckUnique checks if any value is unique in a table
func CheckUnique(db *gorm.DB, column string, value any, tableName string) bool {
	var count int64

	// Counting the occurrences of the value in the specified column
	if err := db.Table(tableName).Where(column+" = ?", value).Count(&count).Error; err != nil {
		log.Panicf("failed to check if %s is unique: %v", column, err)
		return false
	}

	// If count is 0, the value is unique
	return count == 0
}
