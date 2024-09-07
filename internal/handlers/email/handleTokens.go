package email

import (
	"errors"
	"log"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers"
	"gorm.io/gorm"
)

func HandleRegisterTokens(db *gorm.DB, userId int64) (error, string) {
	// generate a new token
	confirmationToken := handlers.GenerateToken()

	if confirmationToken == "" {
		log.Printf("\n\nERROR\n\tFailed to generate token for user: %v\n\n", userId)
		return errors.New("Failed to generate token for user: "), ""
	}

	verificationToken := &handlers.VerificationTokens{
		UserId: userId,
		Token:  confirmationToken,
	}

	// save confirmationToken to the database
	if err := db.Create(&verificationToken).Error; err != nil {
		log.Println("ERROR\n\tFailed to save verification token: ", err)
		return errors.New("Failed to create/save verification token: " + err.Error()), ""
	}

	return nil, confirmationToken
}

// HandleLoginCodes generates a 6 digit code used for email auth
func HandleLoginCodes(db *gorm.DB, userId int64) (error, string) {
	loginCode, err := handlers.GenerateCode()

	if err != nil {
		log.Printf("\n\nERROR\n\tFailed to generate login code for user: %v\n\n", userId)
		return errors.New("Failed to generate login code for user: "), ""
	}

	userCode := &handlers.LoginCodes{
		UserId: userId,
		Code:   loginCode,
	}

	// save loginCode to the database
	if err := db.Create(&userCode).Error; err != nil {
		log.Println("ERROR\n\tFailed to save login code: ", err)
		return errors.New("Failed to create/save login code: " + err.Error()), ""
	}

	return nil, loginCode
}
