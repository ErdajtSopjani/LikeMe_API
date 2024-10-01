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

// GetUserMessagesResponse is the structure of the response
type GetUserMessagesResponse struct {
	User1_ID int64              `json:"user1_id"`
	User2_ID int64              `json:"user2_id"`
	Messages []handlers.Message `json:"messages"`
}

// GetUserMessages returns all the messages between the user and another specific user
func GetUserMessages(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userToken := r.Header.Get("Authorization")
		var req GetUserMessagesRequest // decode the request body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helpers.RespondError(w, "Invalid Request-Format", http.StatusBadRequest)
			log.Printf("\n\nBAD REQUEST\n\tBad request: %v\n\tError: %s\n\n", req, err)
			return
		}

		// get both user Ids
		var user handlers.UserToken
		var messagedUser handlers.UserProfile
		err := db.Where("token = ?", userToken).First(&user).Error
		err = db.Where("username = ?", req.MessagedUsername).First(&messagedUser).Error

		if err != nil {
			helpers.RespondError(w, "User not found", http.StatusBadRequest)
			log.Printf("\n\nBAD/MALICIOUS\n\tBad credential request from %s: %v\n\tError: %s\n\n", r.RemoteAddr, req, err)
			return
		}

		// get all messages between the two users
		var messages []handlers.Message
		err = db.Where("sender_id = ? AND receiver_id  = ?", user.UserId, messagedUser.UserId).Or("sender_id = ? AND receiver_id= ?", messagedUser.UserId, user.UserId).Find(&messages).Order("created_at").Error
		if err != nil {
			helpers.RespondError(w, "Failed to get messages", http.StatusInternalServerError)
			log.Printf("\n\nERROR\n\tFailed to get messages: %v\n\tError: %s\n\n", req, err)
			return
		}

		// Respond with the messages
		response := GetUserMessagesResponse{
			User1_ID: user.UserId,
			User2_ID: messagedUser.UserId,
			Messages: messages,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
