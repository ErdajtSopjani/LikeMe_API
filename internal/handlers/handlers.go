package handlers

import "time"

/* NOTE: TABLE SCHEMA */
// all these structs are used to define the schema of the database tables
// they all represend a table in the database
// i use gorm as an ORM to interact with the database

var Now time.Time = time.Now().UTC()
var VerificationExpiresAt time.Time = Now.Add(10 * time.Minute)
var TokenExpiresAt time.Time = Now.Add(24 * time.Hour)

// User is used to store generic information upon registration
type User struct {
	// ID is bigserial and is the primary key
	ID          int64      `gorm:"type:bigserial;primaryKey"`
	Email       string     `gorm:"type:varchar(255);unique;not null"`
	CountryCode string     `gorm:"type:varchar(5);not null"`
	Verified    bool       `gorm:"not null;default:false"`
	CreatedAt   *time.Time `gorm:"type:timestamp;not null"`
}

// UserProfile is only created once the email is verified and the user is logged in
type UserProfile struct {
	ID             int64  `gorm:"primaryKey"`
	Username       string `gorm:"type:varchar(30);unique;not null"`
	FirstName      string `gorm:"type:varchar(100);not null"`
	LastName       string `gorm:"type:varchar(100);not null"`
	ProfilePicture string
	Bio            string     `gorm:"type:varchar(255)"`
	UserId         int64      `gorm:"not null"`
	CreatedAt      *time.Time `gorm:"type:timestamp;not null"`
}

// UserToken is created to manage user sessions since we only use email authentication
type UserToken struct {
	ID        int64      `gorm:"primaryKey"`
	Token     string     `gorm:"type:varchar(32);unique;not null"`
	UserId    int64      `gorm:"not null"`
	CreatedAt *time.Time `gorm:"type:timestamp;not null"`
	ExpiresAt *time.Time `gorm:"type:timestamp;not null;"`
}

// EmailChange is used to store the email change request and the token for later confirmation
type EmailChange struct {
	ID          int64      `gorm:"primaryKey"`
	UserId      int64      `gorm:"not null"`
	Email       string     `gorm:"type:varchar(255);unique;not null"`
	ChangeToken string     `gorm:"type:varchar(32);unique;not null"`
	CreatedAt   *time.Time `gorm:"type:timestamp;not null"`
}

// Tag is used to categorize posts and for easier implementation of the algorithm
type Tag struct {
	ID   int64  `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(255);unique;not null"`
}

// UserInterest stores user interests based on how they interact with certain Tags
type UserInterest struct {
	ID         int64      `gorm:"primaryKey"`
	UserId     int64      `gorm:"not null"`
	InterestId int64      `gorm:"not null"`
	CreatedAt  *time.Time `gorm:"type:timestamp;not null"`
}

// Follow is used to store the relationship between users
type Follow struct {
	ID          int64      `gorm:"primaryKey"`
	FollowerId  int64      `gorm:"not null"`
	FollowingId int64      `gorm:"not null"`
	CreatedAt   *time.Time `gorm:"type:timestamp;not null"`
}

// BlockedUser keeps track of whether a user is blocked or not and will be used to limit blocked interactions
type BlockedUser struct {
	ID            int64      `gorm:"primaryKey"`
	UserId        int64      `gorm:"not null"`
	BlockedUserId int64      `gorm:"not null"`
	CreatedAt     *time.Time `gorm:"type:timestamp;not null"`
}

// Post is a generic social media post. Supports all media and must contain a tag
type Post struct {
	ID        int64  `gorm:"primaryKey"`
	UserId    int64  `gorm:"not null"`
	Title     string `gorm:"type:varchar(255);not null"`
	Content   string `gorm:"not null"`
	Media     string
	MediaType string     `gorm:"type:varchar(10)"`
	TagId     int64      `gorm:"not null"`
	CreatedAt *time.Time `gorm:"type:timestamp;not null"`
}

// Like is used to store the relationship between a user and a post
type Like struct {
	ID        int64      `gorm:"primaryKey"`
	UserId    int64      `gorm:"not null"`
	PostId    int64      `gorm:"not null"`
	CreatedAt *time.Time `gorm:"type:timestamp;not null"`
}

// Comment is used to store comments on posts; a comment must be linked to a single post
type Comment struct {
	ID        int64      `gorm:"primaryKey"`
	UserId    int64      `gorm:"not null"`
	PostId    int64      `gorm:"not null"`
	Content   string     `gorm:"not null"`
	CreatedAt *time.Time `gorm:"type:timestamp;not null"`
}

// Message is used to store messages between users and will be used for the chat feature
type Message struct {
	ID         int64      `gorm:"primaryKey"`
	SenderId   int64      `gorm:"not null"`
	ReceiverId int64      `gorm:"not null"`
	Content    string     `gorm:"not null"`
	CreatedAt  *time.Time `gorm:"type:timestamp;not null"`
}

// TwoFactor is used to temporarily store 6-digit codes for users to login
type TwoFactor struct {
	ID        int64      `gorm:"primaryKey"`
	Code      int        `gorm:"not null;unique"`
	UserId    int64      `gorm:"not null"`
	CreatedAt *time.Time `gorm:"type:timestamp;not null"`
	ExpiresAt *time.Time `gorm:"type:timestamp;not null"`
}

// VerificationToken is used to store tokens that will be used to verify a user's email
type VerificationToken struct {
	ID        int64      `gorm:"primaryKey"`
	Token     string     `gorm:"type:varchar(32);unique;not null"`
	UserId    int64      `gorm:"not null"`
	CreatedAt *time.Time `gorm:"type:timestamp;not null"`
	ExpiresAt *time.Time `gorm:"type:timestamp;not null"`
}
