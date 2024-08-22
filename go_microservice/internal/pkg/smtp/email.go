package smtp

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

func Send(to string, token string) error {
	from := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	port := os.Getenv("HTTP_PORT")
	fmt.Printf("Email: %v", to)
	verificationURL := fmt.Sprintf("http://localhost:%s/verify?token=%s", port, token)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	subject := "Subject: Email Verification\r\n"
	body := fmt.Sprintf("Please use the following link to verify your email: %s\r\n", verificationURL)
	mime := "MIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\n\r\n"

	msg := []byte(subject + mime + body)

	log.Printf("Attempting to send email to %s", to)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		return err
	}

	return nil
}
