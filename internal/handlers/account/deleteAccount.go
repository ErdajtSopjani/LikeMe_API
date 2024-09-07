package account

import (
	"net/http"

	"gorm.io/gorm"
)

// DeleteAccount delete's all data associated with an account
func DeleteAccount(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not implemented"))
	}
}
