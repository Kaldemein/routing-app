package main

import (
	// "encoding/json"
	"fmt"
	"os"

	// "log"
	// "time"

	// amqp "github.com/rabbitmq/amqp091-go"

	"database/sql"

	_ "github.com/lib/pq"
)

func main() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

}

// func main() {
// 	var conn *amqp.Connection
// 	var err error

// 	// fmt.Println("Successfully connected to PostgreSQL!")

// 	// Подключение к RabbitMQ
// 	for i := 0; i < 10; i++ {
// 		conn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
// 		if err == nil {
// 			log.Printf("Succsessfully connected to RabbitMQ")
// 			break
// 		}
// 		log.Printf("Failed to connect to RabbitMQ: %v", err)
// 		time.Sleep(2 * time.Second)
// 	}
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	if err != nil {
// 		log.Fatalf("Failed to open a channel: %v", err)
// 	}
// 	defer ch.Close()

// 	q, err := ch.QueueDeclare(
// 		"email_verification", // Имя очереди
// 		false,                // durable
// 		false,                // delete when unused
// 		false,                // exclusive
// 		false,                // no-wait
// 		nil,                  // arguments
// 	)
// 	if err != nil {
// 		log.Fatalf("Failed to declare a queue: %v", err)
// 	}

// 	msgs, err := ch.Consume(
// 		q.Name, // queue
// 		"",     // consumer
// 		true,   // auto-ack
// 		false,  // exclusive
// 		false,  // no-local
// 		false,  // no-wait
// 		nil,    // args
// 	)
// 	if err != nil {
// 		log.Fatalf("Failed to register a consumer: %v", err)
// 	}

// 	// Обработка сообщений
// 	forever := make(chan bool)

// 	go func() {
// 		for d := range msgs {
// 			// Обработка сообщения
// 			processEmailVerification(d.Body)

// 			log.Printf("Received a message: %s", d.Body)
// 		}
// 	}()

// 	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
// 	<-forever
// }

// func processEmailVerification(body []byte) {
// 	// Разбор JSON сообщения
// 	var data map[string]string
// 	if err := json.Unmarshal(body, &data); err != nil {
// 		log.Printf("Error decoding JSON: %v", err)
// 		return
// 	}

// 	// Генерация токена и отправка email
// 	email := data["email"]
// 	token := generateToken()

// 	// Сохранение токена в базе данных микросервиса

// 	// Отправка email с ссылкой для подтверждения
// 	sendConfirmationEmail(email, token)
// }

// func generateToken() string {
// 	// Реализация генерации токена
// 	return "example_token"
// }

// func sendConfirmationEmail(email, token string) {
// 	// Реализация отправки email с ссылкой для подтверждения
// 	fmt.Printf("Sending confirmation email to %s with token %s\n", email, token)
// }
