package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func Send(to *string, subject string, text *string) error {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Couldn't load .env variables %v\n", err)
	}
	email := os.Getenv("EMAIL")
	emailAppPassword := os.Getenv("EMAIL_APP_PASSWORD")
	message := fmt.Sprintf("Subject: %s\n %s", subject, *text)

	mailAuth := smtp.PlainAuth(
		"",
		email,
		emailAppPassword,
		"smtp.gmail.com",
	)
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		mailAuth,
		email,
		[]string{*to},
		[]byte(message),
	)
	if err != nil {
		return err
	}
	return nil
}