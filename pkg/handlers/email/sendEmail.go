package email

import (
	"fmt"
	"log"
	"os"

	"github.com/ErdajtSopjani/LikeMe_API/pkg/handlers"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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
func SendConfirmation(userEmail string) {
	var verifyEmail string = fmt.Sprintf("<html><body><h1>Verify your email address</h1><br/><h3>Thank you for becoming part of LikeMe!<h3><p>If you encounter any problems feel free to reach out via this email address.</p><br/><p>Please verify your email address to proceed</p> <a href=\"https://%s/verify?token=%s\">Click here to verify</a><br /> <br />", os.Getenv("FrontEnd_URL"), handlers.GenerateToken())

	fmt.Println(verifyEmail)

	// Create a new email
	confirmationEmail := &Email{
		From:             mail.NewEmail("LikeMe", "verify.likeme.dev@gmail.com"),
		Subject:          "Welcome to LikeMe",
		To:               mail.NewEmail("User", userEmail),
		PlainTextContent: "Please verify your email address",
		HTMLContent:      verifyEmail,
		client:           sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY")),
	}

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
	} else {
		log.Println(response.StatusCode, ": Email confirmation sent: %v", confirmationEmail.To)
	}
}
