package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	apiHandler "service/internal/api"
	"service/internal/pkg/db"
	"service/internal/pkg/rabbitmq"
	"service/internal/pkg/smtp"
	emailtokens "service/internal/repository/email_tokens"
	emailverification "service/internal/usecases/email_verification"
)

func main() {
	// Connection to postgres
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Database connection
	dbConn, err := db.ConnectToDB(psqlInfo)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %s", err)
	}
	defer dbConn.Close()
	log.Print("Postgres connected")

	// RabbitMQ connection
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	rabbitConn, err := rabbitmq.GetConnection(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer rabbitConn.Close()
	log.Print("RabbitMQ connected")

	// Initialize interfaces
	emailRepo := emailtokens.NewPostgresEmailRepository(dbConn)
	rabbitMQ, err := rabbitmq.NewQueue(rabbitConn)
	if err != nil {
		log.Fatalf("Failed to initialize rabbitmq interface: %v", err)
	}

	// Use case
	emailUseCase := emailverification.NewEmailUseCase(emailRepo, rabbitMQ, smtp.Send)

	// HTTP handler
	httpHandler := apiHandler.NewHttpHandler(emailUseCase)

	//Declare queue
	err = rabbitMQ.DeclareQueue("email_queue")
	if err != nil {
		log.Printf("Failed to declare queue: %v", err)
	}

	// Start HTTP server
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "9999"
	}

	go func() {
		log.Printf("HTTP Server started on port %s", httpPort)
		log.Fatal(http.ListenAndServe(":"+httpPort, httpHandler))
	}()

	msgs, err := rabbitMQ.Consume("email_queue")
	if err != nil {
		log.Fatalf("Failed to register a consumer email_queue: %s", err)
	}

	for msg := range msgs {
		msg := msg
		go func() {
			email := string(msg.Body)
			log.Printf("Processing email: %s", email)
			err := emailUseCase.GenerateAndSendLink(email)
			if err != nil {
				log.Printf("Failed to process email %s: %v", email, err)
			}
		}()
	}

	select {}

}
