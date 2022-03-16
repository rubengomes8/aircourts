package smtp

import (
	"crypto/tls"
	"fmt"
	"strings"

	gomail "gopkg.in/mail.v2"
)

type Sender struct {
	Email    string
	Password string
}

type Email struct {
	To      string
	From    string
	Subject string
	Body    string
}

func NewSender(email, password string) Sender {
	return Sender{
		Email:    email,
		Password: password,
	}
}

func (s Sender) SendEmail(email Email) error {

	if !Valid(email.From) || !Valid(email.To) {
		return fmt.Errorf("From %v or To %v invalid", email.From, email.To)
	}

	m := gomail.NewMessage()

	m.SetHeader("From", email.From)
	m.SetHeader("To", email.To)
	m.SetHeader("Subject", email.Subject)
	m.SetBody("text/plain", email.Body)

	d := gomail.NewDialer("smtp.gmail.com", 587, s.Email, s.Password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("error dialing and sending email: %v", err.Error())
	}

	return nil
}

func Valid(email string) bool {
	if email == "" {
		return false
	}

	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}

	return true
}
