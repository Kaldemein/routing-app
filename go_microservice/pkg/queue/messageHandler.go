package queue

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"log"
	"service/pkg/emails"
	"service/pkg/errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func MessageHandler(db *sql.DB, msg amqp.Delivery) {
	email := msg.Body
	log.Printf("Received a message: %s", email)

	token := generateSecureToken(20)
	log.Printf("Generated token:: %s", token)

	log.Printf("checking email if exists in db...")
	// Check if the email already exists in the database
	var exists bool
	queryCheck := `SELECT EXISTS (SELECT 1 FROM email_tokens WHERE email=$1)`
	err := db.QueryRow(queryCheck, email).Scan(&exists)
	if err != nil {
		log.Printf("Failed to check email existence: %s", err)
		return
	}
	if exists {
		log.Printf("Email already exists: %s", email)
	} else {
		//saving email, token to DB
		query := `INSERT INTO email_tokens (email, token) VALUES ($1, $2)`
		_, err = db.Exec(query, email, token)
		errors.FailOnError(err, "Failed insert data :(")
		email := string(msg.Body)
		err = emails.SendEmail(email, token)
		errors.FailOnError(err, "Failed send email :(")
		log.Printf("check your email!!!")
	}
}
