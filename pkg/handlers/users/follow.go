package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ErdajtSopjani/LikeMe_API/pkg/handlers"
	"gorm.io/gorm"
)

// FollowAccount is a handler that checks if both users exist and creates a new follow in the database
func FollowAccount(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req handlers.Follows

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil { // if the request body is not valid
			log.Fatal("FollowAccount failed to decode request body: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// check if both users exist by checking if their values are not unique
		usersExist := !handlers.CheckUnique(db, "id", req.FollowerId, "users") && !handlers.CheckUnique(db, "id", req.FollowingId, "users")
		if !usersExist { // if one or both values return to be unique
			http.Error(w, "Invalid user_id", http.StatusBadRequest)
			return
		}

		follows := handlers.Follows{
			FollowerId:  req.FollowerId,
			FollowingId: req.FollowingId,
		}

		if err := db.Create(&follows).Error; err != nil { // create a new follow
			log.Fatal("failed to create new follow: ", err)
			http.Error(w, "Internal-Server Error...", http.StatusInternalServerError)
			return
		}

		// return status created and the new follow
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(w); err != nil {
			log.Fatal("failed to encode response: ", err)
			http.Error(w, "Internal-Server Error...", http.StatusInternalServerError)
			return
		}
	}
}
