package email

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ErdajtSopjani/LikeMe_API/pkg/handlers"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"gorm.io/gorm"
)

// Email is the structure accepted by sendgrid to send an email
type Email struct {
	From             *mail.Email
	Subject          string
	To               *mail.Email
	PlainTextContent string
	HTMLContent      string
	client           *sendgrid.Client
}

// SendConfirmation sends an email verification on account creation
func SendConfirmation(db *gorm.DB, userEmail string, userId int64) error {
	err, confirmationToken := handleTokens(db, userId)
	if err != nil {
		return err
	}

	var verifyEmail string = fmt.Sprintf(`
        <html>
            <body>
                <h1>Verify your email address</h1>
                <br/>
                <h3>Thank you for becoming part of LikeMe!<h3>

                <br/>
                <p>Please verify your email address to proceed</p>
                <a href="%s/verify?token=%s">Click here to verify</a>

                <br/>
                <br/>
                <p>If you encounter any problems feel free to reach out via this email address.</p>
                <br/>
                <br/>
                <p>All the best,<br/>LikeMe</p>
                <br/>
            </body>
        </html>`, os.Getenv("FRONTEND_URL"), confirmationToken)

	// Create a new email
	confirmationEmail := &Email{
		From:             mail.NewEmail("LikeMe", "verify.likeme.dev@gmail.com"),
		Subject:          "Welcome to LikeMe",
		To:               mail.NewEmail("User", userEmail),
		PlainTextContent: "Please verify your email address",
		HTMLContent:      verifyEmail,
		client:           sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY")),
	}

	// TODO: Check if the email and user_id match before sending the email

	// Send the email
	response, err := confirmationEmail.client.Send(
		mail.NewSingleEmail(
			confirmationEmail.From,
			confirmationEmail.Subject,
			confirmationEmail.To,
			confirmationEmail.PlainTextContent,
			confirmationEmail.HTMLContent),
	)

	if err != nil {
		log.Printf("\n\nERROR\n\tUnable to confirmation email\n\tERROR: %v\n\n", err)
		return err
	} else {
		log.Println(response.StatusCode, ": Email confirmation sent: %v", confirmationEmail.To)
	}

	return nil
}

// handleTokens generates a new token and saves it to the database
func handleTokens(db *gorm.DB, userId int64) (error, string) {
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
