package verify

import (
	"fmt"
	"log"
	"os"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/email"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"gorm.io/gorm"
)

// SendConfirmation sends an email verification on account creation
func SendConfirmation(db *gorm.DB, userEmail string, userId int64) error {
	err, confirmationToken := email.HandleRegisterTokens(db, userId)
	if err != nil {
		return err
	}

	var verifyEmail string = fmt.Sprintf(`
<html>
    <body style="font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px; color: #333;">
        <div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 10px; box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); padding: 20px;">
            <h1 style="color: #4CAF50; text-align: center; font-size: 28px; margin-bottom: 20px;">Verify Your Email Address</h1>
            <p style="font-size: 16px; color: #555; text-align: center;">Welcome to LikeMe!</p>
            <p style="font-size: 16px; color: #555; text-align: center;">
                Thank you for joining our community. To get started, please verify your email address by clicking the button below.
            </p>
            <div style="text-align: center; margin: 30px 0;">
                <a href="%s/verify?token=%s" style="background-color: #4CAF50; color: #ffffff; text-decoration: none; padding: 15px 30px; font-size: 16px; border-radius: 5px; display: inline-block;">
                    Verify Email
                </a>
            </div>
            <p style="font-size: 14px; color: #777; text-align: center;">
                If the button doesn't work, copy and paste the following URL into your browser:
            </p>
            <p style="font-size: 14px; color: #4CAF50; word-wrap: break-word; text-align: center;">%s/verify?token=%s</p>
            <p style="font-size: 14px; color: #777; text-align: center;">
                If you encounter any issues, feel free to reply to this email for support.
            </p>
            <br />
            <p style="font-size: 14px; color: #777; text-align: center;">
                Best regards,<br />
                <strong>LikeMe Team</strong>
            </p>
        </div>
    </body>
</html>
`, os.Getenv("FRONTEND_URL"), confirmationToken, os.Getenv("FRONTEND_URL"), confirmationToken)

	// Create a new email
	confirmationEmail := email.Email{
		From:             mail.NewEmail("LikeMe", "verify.likeme.dev@gmail.com"),
		Subject:          "Welcome to LikeMe",
		To:               mail.NewEmail("User", userEmail),
		PlainTextContent: "Please verify your email address",
		HTMLContent:      verifyEmail,
		Client:           sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY")),
	}

	// check if the email and user_id match before sending the email
	// this is more of a fail-safe since we get the userId directly from the database
	if err := db.Where("email = ?", userEmail).First(&handlers.User{}).Error; err != nil {
		log.Printf("\n\nERROR\n\tFailed to find user with email: %v\n\tError: %v\n\n", userEmail, err)
		return err
	}

	// Send the email
	response, err := confirmationEmail.Client.Send(
		mail.NewSingleEmail(
			confirmationEmail.From,
			confirmationEmail.Subject,
			confirmationEmail.To,
			confirmationEmail.PlainTextContent,
			confirmationEmail.HTMLContent),
	)

	if err != nil || response.StatusCode != 202 {
		log.Printf("\n\nERROR\n\tUnable to send confirmation email\n\tERROR: %v\n\tSTATUS CODE: %v\n\n", err, response.StatusCode)
		return err
	} else {
		log.Println(response.StatusCode, ": Email confirmation sent: %v", confirmationEmail.To)
	}

	return nil
}
