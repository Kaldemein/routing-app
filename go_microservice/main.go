package main

import (
	// "encoding/json"

	"fmt"
	"log"
	"net/http"
	"os"
	"service/pkg/database"
	"service/pkg/emails"
	"service/pkg/errors"
	"service/pkg/queue"

	"time"

	_ "github.com/lib/pq"
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

	db, err := database.ConnectToDB(psqlInfo, retryCount, retryInterval)
	errors.FailOnError(err, "Failed to connect postgres :(")
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
	errors.FailOnError(err, "Failed to create table")
	fmt.Println("Table created: OK!")

	//RabbitMQ connection
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	// Set retry parameters for rabbitmq connection
	retryCount = 5
	retryInterval = 5 * time.Second

	rabbitMQConn, err := queue.ConnectToRabbitMQ(rabbitMQURL, retryCount, retryInterval)
	errors.FailOnError(err, "Failed to connect to RabbitMQ")
	defer rabbitMQConn.Close()
	fmt.Println("RabbitMQ: OK!")

	// Create a channel
	rabbitMQCh, err := rabbitMQConn.Channel()
	errors.FailOnError(err, "Failed to open a channel")
	defer rabbitMQCh.Close()

	// Declare a queue
	queueName := "email_queue"
	_, err = rabbitMQCh.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	errors.FailOnError(err, "Failed to declare a queue")

	// Consume messages from the queue
	msgs, err := rabbitMQCh.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	errors.FailOnError(err, "Failed to register a consumer")

	// email verification endpoint
	http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		emails.VerifyHandler(db, w, r, rabbitMQCh)
		log.Printf("EMAIL VERIFIED")
	})

	// http server settings
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "9999"
	}

	// start http server
	go func() {
		log.Printf("HTTP Server started on port %s", httpPort)
		log.Fatal(http.ListenAndServe(":"+httpPort, nil))
	}()

	//take message from queue and handle it
	for msg := range msgs {
		log.Printf("in msg cycle")
		go queue.MessageHandler(db, msg)
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
}
