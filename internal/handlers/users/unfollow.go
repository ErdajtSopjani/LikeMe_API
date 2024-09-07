package users

import (
	"net/http"

	"gorm.io/gorm"
)

func UnfollowAccount(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
