package helpers

import "gorm.io/gorm"

// UsersExist just checks if given userIds are in the database
func UsersExist(userIds []int64, db *gorm.DB) bool {
	for _, id := range userIds {
		var count int64
		db.Table("users").Where("id = ?", id).Count(&count)
		if count != 1 {
			return false
		}
	}

	return true
}
