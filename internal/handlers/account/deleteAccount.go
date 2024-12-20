package account

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/helpers"
	"gorm.io/gorm"
)

type DeleteRequest struct {
	Email string `json:"email"`
}

// DeleteAccount delete's all data associated with a user's id and the user itself
func DeleteAccount(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		var req DeleteRequest
		// decode the request body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helpers.RespondError(w, "Invalid Request-Format", http.StatusBadRequest) // if the request body is not valid
			log.Printf("\n\nBAD REQUEST\n\tBad request: %v\n\tError: %s\n\n", req, err)
			return
		}

		// check if the email exists in the database
		if helpers.CheckUnique(db, "email", req.Email, "users") {
			helpers.RespondError(w, "Email not found", http.StatusBadRequest)
			log.Printf("\n\nBAD/MALICIOUS\n\tBad credential request from %s: %v\n\tError: %s\n\n", r.RemoteAddr, req, "Email not found")
			return
		}

		// get the userId from the email
		var user handlers.User
		if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
			helpers.RespondError(w, "User not found", http.StatusBadRequest)
			log.Printf("\n\nBAD/MALICIOUS\n\tBad credential request from %s: %v\n\tError: %s\n\n", req, err)
		}

		// check if the token with the user id exists
		if err := db.Where("token = ? AND user_id = ?", token, user.ID).First(&handlers.UserToken{}).Error; err != nil {
			helpers.RespondError(w, "Invalid Token", http.StatusBadRequest)
			log.Printf("\n\nBAD/MALICIOUS\n\tBad credential request from %s: %v\n\tError: %s\n\n", r.RemoteAddr, req, "Invalid Token")
			return
		}

		err := deleteAllData(db, w, req, user)
		if err != nil {
			helpers.RespondError(w, "Failed to delete user data", http.StatusInternalServerError)
			log.Printf("\n\nERROR\n\tFailed to delete user data: %v\n\tError: %s\n\n", req, err)
			return
		}

		// Respond with success
		helpers.RespondJSON(w, http.StatusOK, "User deleted successfully")
	}
}

// deleteAllData deletes all data associated to a user in all columns
func deleteAllData(db *gorm.DB, w http.ResponseWriter, req DeleteRequest, user handlers.User) error {
	// begin a db transaction
	tx := db.Begin()

	tablesToDeleteFrom := []string{
		"user_profiles", "user_tokens", "user_interests", "follows",
		"blocked_users", "posts", "likes", "comments", "messages",
		"two_factors", "verification_tokens",
	}

	for _, table := range tablesToDeleteFrom {
		if table == "follows" {
			// Delete records where the user is the follower or the following
			if err := tx.Table(table).Where("follower_id = ? OR following_id = ?", user.ID, user.ID).Delete(nil).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Rollback()
				log.Printf("Error deleting from table %s for user_id %d: %s", table, user.ID, err)
				return err
			}
		} else if table == "messages" {
			// Delete messages where the user is the sender or receiver
			if err := tx.Table(table).Where("sender_id = ? OR receiver_id = ?", user.ID, user.ID).Delete(nil).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Rollback()
				log.Printf("Error deleting from table %s for user_id %d: %s", table, user.ID, err)
				return err
			}
		} else {
			// Default deletion logic
			if err := tx.Table(table).Where("user_id = ?", user.ID).Delete(nil).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Rollback()
				log.Printf("Error deleting from table %s for user_id %d: %s", table, user.ID, err)
				return err
			}
		}
	}

	// delete the user from the users table
	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		log.Printf("Error deleting user %d: %s", user.ID, err)
		return err
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Printf("Transaction failed for user %d: %s", user.ID, err)
		return err
	}

	return nil
}
