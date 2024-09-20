package users

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/helpers"
	"gorm.io/gorm"
)

type FollowRequest struct {
	FollowerId  int64 `json:"follower_id"`
	FollowingId int64 `json:"following_id"`
}

// FollowAccount is a handler that checks if both users exist and creates a new follow in the database
func FollowAccount(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req FollowRequest

		// decode request body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helpers.RespondError(w, "Invalid request format", http.StatusBadRequest)
			log.Printf("\n\nBAD REQUEST\n\tBad request on follow: %v\n\tError: %s\n\n", req, err)
			return
		}

		// check if both users exist
		var count int64
		db.Table("users").Where("id IN (?, ?)", req.FollowerId, req.FollowingId).Count(&count)
		if count != 2 {
			helpers.RespondError(w, "Invalid follower or following ID", http.StatusBadRequest)
			return
		}

		// check if token and userId match
		if !helpers.CheckTokenMatch(db, req.FollowerId, r.Header.Get("Authorization")) {
			helpers.RespondError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// check if follow record already exists
		if err := db.Where("follower_id = ? AND following_id = ?", req.FollowerId, req.FollowingId).First(&handlers.Follow{}).Error; err == nil {
			helpers.RespondError(w, "Follow already exists", http.StatusBadRequest)
			return
		}

		// create new follow record and save it to the database
		follow := handlers.Follow{
			FollowerId:  req.FollowerId,
			FollowingId: req.FollowingId,
		}
		if err := db.Create(&follow).Error; err != nil {
			helpers.RespondError(w, "Failed to create follow", http.StatusInternalServerError)
			log.Printf("\n\nERROR\n\tFailed to create follow: %v\n\tError: %s\n\n", follow, err)
			return
		}

		// respond with status created and the new follow record
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(follow); err != nil {
			log.Printf("\n\nERROR\n\tFailed to encode follow: %v\n\tError: %s\n\n", follow, err)
			helpers.RespondError(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
