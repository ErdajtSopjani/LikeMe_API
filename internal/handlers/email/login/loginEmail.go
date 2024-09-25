package login_email

import (
	"fmt"
	"log"
	"os"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/email"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"gorm.io/gorm"
)

// SendLoginEmail sends an email for logging in
func SendLoginEmail(db *gorm.DB, userId int64, userEmail string) error {
	err, loginCode := email.HandleLoginCodes(db, userId)
	if err != nil {
		return err
	}

	var loginEmailHTML string = fmt.Sprintf(`<html>

<body style="font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px; color: #333;">
    <div
        style="max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 10px; box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); padding: 20px;">
        <h1 style="color: #4CAF50; text-align: center; font-size: 28px; margin-bottom: 20px;">Finish Logging In</h1>
        <p style="font-size: 16px; color: #555; text-align: center;">Hi there,</p>
        <p style="font-size: 16px; color: #555; text-align: center;">
            You’re just one step away from accessing your account. Click the button below to log in or use the code:
            <strong>%d</strong>.
        </p>
        <div style="text-align: center; margin: 30px 0;">
            <a href="%s/verify?code=%d"
                style="background-color: #4CAF50; color: #ffffff; text-decoration: none; padding: 15px 30px; font-size: 16px; border-radius: 5px; display: inline-block;">
                Click Here to Log In
            </a>
        </div>
        <p style="font-size: 14px; color: #777; text-align: center;">
            If the button doesn't work, copy and paste the following URL into your browser:
        </p>
        <p style="font-size: 14px; color: #4CAF50; word-wrap: break-word; text-align: center;">%s/verify?code=%d</p>
        <p style="font-size: 14px; color: #777; text-align: center;">
            If you encounter any issues, feel free to reply to this email for support.
        </p>
        <br />
        <p style="font-size: 14px; color: #777; text-align: center;">
            All the best,<br />
            <strong>LikeMe</strong>
        </p>
    </div>
</body>

</html>`, loginCode, os.Getenv("FRONTEND_URL"), loginCode, os.Getenv("FRONTEND_URL"), loginCode)

	loginEmail := &email.Email{
		From:             mail.NewEmail("LikeMe", "verify.likeme.dev@gmail.com"),
		Subject:          fmt.Sprintf("%d is your code to log in to LikeMe", loginCode),
		To:               mail.NewEmail("User", userEmail),
		PlainTextContent: "Use this code to log in",
		HTMLContent:      loginEmailHTML,
		Client:           sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY")),
	}

	response, err := loginEmail.Client.Send(
		mail.NewSingleEmail(
			loginEmail.From,
			loginEmail.Subject,
			loginEmail.To,
			loginEmail.PlainTextContent,
			loginEmail.HTMLContent),
	)

	if err != nil {
		log.Printf("\n\nERROR\n\tUnable to send confirmation email\n\tERROR: %v\n\n", err)
		return err
	} else {
		log.Println(response.StatusCode, ": Login Code confirmation sent: %v", loginEmail.To)
	}

	return nil
}
