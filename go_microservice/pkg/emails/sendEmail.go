package emails

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SendEmail(to string, token string) error {
	log.Printf("SEND EMAIL FUNCTION 1")
	from := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	port := os.Getenv("HTTP_PORT")

	verificationURL := fmt.Sprintf("http://localhost:%s/verify?token=%s", port, token)
	log.Printf("verificationURL: %s", verificationURL)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	subject := "Subject: Email Verification\n"
	body := fmt.Sprintf("Please use the following link to verify your email: %s\n", verificationURL)

	mime := "MIME-version: 1.0;\r\nContent-Type: text/plain; charset=\"UTF-8\";\r\n\r\n"

	msg := []byte(subject + mime + body)

	log.Printf("Attempting to send email to %s", to)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}
	log.Printf("SEND EMAIL FUNCTION 2")
	return nil
}

func SendVerificationMessage(email string, rabbitMQCh *amqp.Channel) error {
	log.Printf("sendVerificationMessage is working!!!")
	queueName := "verification_queue"
	_, err := rabbitMQCh.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	body := email
	err = rabbitMQCh.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	return err
}
