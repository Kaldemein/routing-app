package emails

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"service/pkg/errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

func VerifyHandler(db *sql.DB, w http.ResponseWriter, r *http.Request, rabbitMQCh *amqp.Channel) {
	if rabbitMQCh == nil {
		http.Error(w, "RabbitMQ channel not initialized", http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	// email checking and verification
	var email string
	query := `SELECT email FROM email_tokens WHERE token=$1`
	err := db.QueryRow(query, token).Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid or expired token", http.StatusBadRequest)
			return
		}
		errors.FailOnError(err, "Failed to query token")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = SendVerificationMessage(email, rabbitMQCh)
	if err != nil {
		log.Printf("Error sending verification message: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Email verified successfully")
}
