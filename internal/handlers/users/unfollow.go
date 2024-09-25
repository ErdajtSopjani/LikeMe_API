package users

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/helpers"
	"gorm.io/gorm"
)

// UnfollowRequest is the structure that this endpoint expects
type UnfollowRequest struct {
	FollowerId  int64 `json:"follower_id"`
	FollowingId int64 `json:"following_id"`
}

// UnfollowAccount deletes a follow record from the database which basically unfollows an account
func UnfollowAccount(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// decode the request body
		var req UnfollowRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helpers.RespondError(w, "Invalid request format", http.StatusBadRequest)
			log.Printf("\n\nBAD REQUEST\n\tBad request on unfollow: %v\n\tError: %s\n\n", req, err)
			return
		}

		// check if both users exist, NOTE: although an impossible scenario this is just a failsafe.
		if helpers.UsersExist([]int64{req.FollowerId, req.FollowingId}, db) == false {
			helpers.RespondError(w, "Invalid follower or following ID", http.StatusBadRequest)
			return
		}

		// check if the token and userId match
		if !helpers.CheckTokenMatch(db, req.FollowerId, r.Header.Get("Authorization")) {
			helpers.RespondError(w, "Unauthorized", http.StatusUnauthorized)
			log.Printf("\n\nUNAUTHORIZED\n\tUnauthorized request on unfollow: %v\n\n", req)
			return
		}

		// check if the follow record exists
		if err := db.Where("follower_id = ? AND following_id = ?", req.FollowerId, req.FollowingId).First(&handlers.Follow{}).Error; err != nil {
			helpers.RespondError(w, "Follow does not exist", http.StatusBadRequest)
			log.Printf("\n\nERROR\n\tFollow does not exist: %v\n\tError: %s\n\n", req, err)
			return
		}

		// delete the follow record
		if err := db.Where("follower_id = ? AND following_id = ?", req.FollowerId, req.FollowingId).Delete(&handlers.Follow{}).Error; err != nil {
			helpers.RespondError(w, "Failed to delete follow", http.StatusInternalServerError)
			log.Printf("\n\nERROR\n\tFailed to delete follow: %v\n\tError: %s\n\n", req, err)
			return
		}

		// Respond with success
		helpers.RespondJSON(w, http.StatusOK, "Unfollowed successfully")
	}
}
