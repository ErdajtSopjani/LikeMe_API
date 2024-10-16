package email

import (
	"errors"
	"log"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/helpers"
	"gorm.io/gorm"
)

// HandleRegisterTokens generates a token for email verification and saves it to the database
func HandleRegisterTokens(db *gorm.DB, userId int64) (string, error) {
	// generate a new token
	confirmationToken := helpers.GenerateToken()
	if confirmationToken == "" {
		log.Printf("\n\nERROR\n\tFailed to generate token for user: %v\n\n", userId)
		return "", errors.New("Failed to generate token for user: ")
	}

	// create and save the token that's available for 10 minutes
	verificationToken := &handlers.VerificationToken{
		UserId:    userId,
		Token:     confirmationToken,
		CreatedAt: &handlers.Now,
		ExpiresAt: &handlers.VerificationExpiresAt,
	}

	// save confirmationToken to the database
	if err := db.Create(&verificationToken).Error; err != nil {
		log.Println("ERROR\n\tFailed to save verification token: ", err)
		return "", errors.New("Failed to create/save verification token: " + err.Error())
	}

	return confirmationToken, nil
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
		CreatedAt: &handlers.Now,
		ExpiresAt: &handlers.VerificationExpiresAt,
	}

	// delete all codes associated with the user
	if err := db.Where("user_id = ?", userId).Delete(&handlers.TwoFactor{}).Error; err != nil {
		log.Println("ERROR\n\tFailed to delete old login codes: ", err)
		return errors.New("Failed to delete old login codes: " + err.Error()), 0
	}

	// save loginCode to the database
	if err := db.Create(&twoFactor).Error; err != nil {
		log.Println("ERROR\n\tFailed to save login code: ", err)
		return errors.New("Failed to create/save login code: " + err.Error()), 0
	}

	return nil, loginCode
}
