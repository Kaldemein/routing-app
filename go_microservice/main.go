package main

import (
	// "encoding/json"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"

	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Connection to postgres
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Set retry parameters for postgres connection
	retryCount := 5
	retryInterval := 5 * time.Second

	db, err := connectToDB(psqlInfo, retryCount, retryInterval)
	failOnError(err, "Failed to connect postgres :(")
	fmt.Println("Postgres: OK!")
	defer db.Close()

	//Create table
	createTableSQL := `CREATE TABLE IF NOT EXISTS email_tokens (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) NOT NULL UNIQUE,
		token VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createTableSQL)
	failOnError(err, "Failed to create table")
	fmt.Println("Table created: OK!")

	//RabbitMQ connection
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	// Set retry parameters for rabbitmq connection
	retryCount = 5
	retryInterval = 5 * time.Second

	rabbitMQConn, err := connectToRabbitMQ(rabbitMQURL, retryCount, retryInterval)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer rabbitMQConn.Close()
	fmt.Println("RabbitMQ: OK!")

	// Create a channel
	ch, err := rabbitMQConn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare a queue
	queueName := "email_queue"
	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Consume messages from the queue
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	failOnError(err, "Failed to register a consumer")

	for msg := range msgs {
		go messageHandler(db, msg)
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
}

func messageHandler(db *sql.DB, msg amqp.Delivery) {
	email := msg.Body
	log.Printf("Received a message: %s", email)

	token := generateSecureToken(20)
	log.Printf("Generated token:: %s", token)
	// Check if the email already exists in the database
	var exists bool
	queryCheck := `SELECT EXISTS (SELECT 1 FROM email_tokens WHERE email=$1)`
	err := db.QueryRow(queryCheck, email).Scan(&exists)
	if err != nil {
		log.Printf("Failed to check email existence: %s", err)
	}
	if exists {
		log.Printf("Email already exists: %s", email)
	} else {
		//saving email, token to DB
		query := `INSERT INTO email_tokens (email, token) VALUES ($1, $2)`
		_, err = db.Exec(query, email, token)
		failOnError(err, "Failed insert data :(")
		email := string(msg.Body) // Преобразование []byte в string
		sendEmail(email, token)
	}
}

func sendEmail(to string, token string) error {
	from := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	subject := "Subject: Email Verification\n"
	body := fmt.Sprintf("Please use the following link to verify your email: %s\n", token)

	mime := "MIME-version: 1.0;\r\nContent-Type: text/plain; charset=\"UTF-8\";\r\n\r\n"

	msg := []byte(subject + mime + body)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		return err
	}

	return nil
}

// Error handler
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// Generate Token
func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func connectToDB(psqlInfo string, retryCount int, retryInterval time.Duration) (*sql.DB, error) {
	var db *sql.DB
	var err error
	for i := 0; i < retryCount; i++ {
		db, err = sql.Open("postgres", psqlInfo)
		if err == nil {
			err = db.Ping()
			if err == nil {
				return db, nil
			}
		}
		fmt.Printf("Failed to connect to database. Retrying in %v seconds...\n", retryInterval.Seconds())
		time.Sleep(retryInterval)
	}
	return nil, err
}

func connectToRabbitMQ(rabbitMQURL string, retryCount int, retryInterval time.Duration) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error
	for i := 0; i < retryCount; i++ {
		conn, err = amqp.Dial(rabbitMQURL)
		if err == nil {
			return conn, nil
		}
		fmt.Printf("Failed to connect to RabbitMQ. Retrying in %v seconds...\n", retryInterval.Seconds())
		time.Sleep(retryInterval)
	}
	return nil, err
}
