package messages

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/helpers"
	"gorm.io/gorm"
)

// GetUserMessagesRequest is the structure GetUserMessages accepts for the request
type GetUserMessagesRequest struct {
	MessagedUsername string `json:"messaged_username"`
}

// GetUserMessages returns all the messages between the user and another specific user
func GetUserMessages(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userToken := r.Header.Get("Authorization")
		var req GetUserMessagesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helpers.RespondError(w, "Invalid Request-Format", http.StatusBadRequest)
			log.Printf("\n\nBAD REQUEST\n\tBad request: %v\n\tError: %s\n\n", req, err)
			return
		}

		// get both user Ids
		var user handlers.User
		var messagedUser handlers.User
		err := db.Where("token = ?", userToken).First(&user).Error
		err = db.Where("username = ?", req.MessagedUsername).First(&messagedUser).Error

		if err != nil {
			helpers.RespondError(w, "User not found", http.StatusBadRequest)
			log.Printf("\n\nBAD/MALICIOUS\n\tBad credential request from %s: %v\n\tError: %s\n\n", r.RemoteAddr, req, err)
			return
		}

		// TODO: get all instances where the both the sender_id and reciever_id is one of the users
	}
}
