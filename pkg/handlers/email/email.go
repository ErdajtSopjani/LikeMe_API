package email

import (
	"log"
	"net/http"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type ConfirmationEmail struct {
	From             *mail.Email
	Subject          string
	To               *mail.Email
	PlainTextContent string
	HTMLContent      string
	message          *mail.Email
	client           *sendgrid.Client
}

// SendConfirmation sends an email verification on account creation
func SendConfirmation(userEmail string) {
	confirmationEmail := &ConfirmationEmail{
		From:             mail.NewEmail("LikeMe", "verify.likeme.dev@gmail.com"),
		Subject:          "Welcome to LikeMe",
		To:               mail.NewEmail("User", userEmail),
		PlainTextContent: "Please verify your email address",
		HTMLContent:      "<strong>and easy to do anywhere, even with Go</strong>",
		client:           sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY")),
	}

	response, err := confirmationEmail.client.Send(
		mail.NewSingleEmail(confirmationEmail.From, confirmationEmail.Subject, confirmationEmail.To, confirmationEmail.PlainTextContent, confirmationEmail.HTMLContent),
	)

	if err != nil {
		log.Fatalf("error: %v", err)
	} else {
		log.Println(response.StatusCode, ": Email confirmation sent: %v", confirmationEmail.To)
	}
}

// VerifyEmail awaits for the SendConfirmation email to be verified
func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

// SendLoginEmail sends an email for logging in
func SendLoginEmail(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}
