package handlers

import "time"

// User represents the structure of the {users} table in the database
type User struct {
	ID          int64  `json:"id"`
	Email       string `json:"email"`
	CountryCode string `json:"country_code"`
	Verified    bool   `json:"is_verified"`
}

// UserProfile represents the structure of user_profile table in the database
type UserProfile struct {
	ID             int64  `json:"id"`
	Username       string `json:"username"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture string `json:"profile_picture"`
	Bio            string `json:"bio"`
	UserId         int64  `json:"user_id"`
}

// Follows represents the follows table in the database
type Follows struct {
	ID          int64 `json:"id"`
	FollowerId  int64 `json:"follower_id"`
	FollowingId int64 `json:"following_id"`
}

// UserTokens represents the verification_tokens and user_tokens tables in the database
type UserTokens struct {
	ID        int64     `json:"id"`
	Token     string    `json:"token"`
	UserId    int64     `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

// VerificationTokens represents the verification_tokens table in the database
type VerificationTokens struct {
	ID        int64     `json:"id"`
	Token     string    `json:"token"`
	UserId    int64     `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

// TwoFactor represents the two_factor table in the database used for email login
type TwoFactor struct {
	ID        int64     `json:"id"`
	Code      string    `json:"code"`
	UserId    int64     `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}
