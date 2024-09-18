package email

import (
	"errors"
	"log"
	"time"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/helpers"
	"gorm.io/gorm"
)

// HandleRegisterTokens generates a token for email verification and saves it to the database
func HandleRegisterTokens(db *gorm.DB, userId int64) (error, string) {
	// generate a new token
	confirmationToken := helpers.GenerateToken()

	if confirmationToken == "" {
		log.Printf("\n\nERROR\n\tFailed to generate token for user: %v\n\n", userId)
		return errors.New("Failed to generate token for user: "), ""
	}

	verificationToken := &handlers.VerificationToken{
		UserId: userId,
		Token:  confirmationToken,
	}

	log.Printf("DEBUG: Attempting to create verification token for userId: %d", userId)

	// save confirmationToken to the database
	if err := db.Create(&verificationToken).Error; err != nil {
		log.Println("ERROR\n\tFailed to save verification token: ", err)
		return errors.New("Failed to create/save verification token: " + err.Error()), ""
	}

	return nil, confirmationToken
}

// HandleLoginCodes generates a 6 digit code used for email auth and saves it to the database
func HandleLoginCodes(db *gorm.DB, userId int64) (error, int) {
	loginCode, err := helpers.GenerateCode()

	if err != nil {
		log.Printf("\n\nERROR\n\tFailed to generate login code for user: %v\n\n", userId)
		return errors.New("Failed to generate login code for user: "), 0
	}

	twoFactor := &handlers.TwoFactor{
		UserId:    userId,
		Code:      loginCode,
		ExpiresAt: time.Now().Add(time.Minute * 10),
	}

	// save loginCode to the database
	if err := db.Create(&twoFactor).Error; err != nil {
		log.Println("ERROR\n\tFailed to save login code: ", err)
		return errors.New("Failed to create/save login code: " + err.Error()), 0
	}

	return nil, loginCode
}
