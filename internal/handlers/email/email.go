package email

import (
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
