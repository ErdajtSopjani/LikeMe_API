package handlers

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

// Follows represents the follows table in the database
type Follows struct {
	FollowerId  int64 `json:"follower_id"`
	FollowingId int64 `json:"following_id"`
}

// VerificationTokens represents the verification_tokens and user_tokens tables in the database
type VerificationTokens struct {
	Token  string `json:"token"`
	UserId int64  `json:"user_id"`
}
