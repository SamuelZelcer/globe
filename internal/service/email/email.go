package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

type Email interface {
	SendVerificationCode(code *string, to *string)
}

type email struct {
	MAIL *string
	MAILAPPPASSWORD *string
}

func InitEmail() Email {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Couldn't load .env variables %v\n", err)
	}
	mail := os.Getenv("EMAIL")
	mailAppPassword := os.Getenv("EMAIL_APP_PASSWORD")
	return &email{
		MAIL: &mail,
		MAILAPPPASSWORD: &mailAppPassword,
	}
}

func (e *email) SendVerificationCode(code *string, to *string) {
	message := fmt.Sprintf("Subject: Verification Code --> \n %s", *code)
	mailAuth := smtp.PlainAuth(
		"",
		*e.MAIL,
		*e.MAILAPPPASSWORD,
		"smtp.gmail.com",
	)
	smtp.SendMail(
		"smtp.gmail.com:587",
		mailAuth,
		*e.MAIL,
		[]string{*to},
		[]byte(message),
	)
}