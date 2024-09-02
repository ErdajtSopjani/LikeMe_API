package handlers

import (
	"time"
)

// User represents the structure of the {users} table in the database
type User struct {
	Email       string `json:"email"`
	CountryCode string `json:"country_code"`
}

// UserProfile represents the structure of user_profile table in the database
type UserProfile struct {
	Username       string `json:"username"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture string `json:"profile_picture"`
	Bio            string `json:"bio"`
	UserId         int64  `json:"user_id"`
}

// UserTokens represents the structure of user_tokens table in the database
type UserTokens struct {
	Id        int64
	ExpiresAt time.Time
	CreatedAt time.Time
	UserId    int64
}

// Follows is the struct for the follows table in the database
type Follows struct {
	FollowerId  int64 `json:"follower_id"`
	FollowingId int64 `json:"following_id"`
}
