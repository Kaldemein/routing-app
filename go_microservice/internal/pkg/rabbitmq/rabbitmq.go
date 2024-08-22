package rabbitmq

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func GetConnection(rabbitMQURL string) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error
	for i := 0; i < 5; i++ {
		conn, err = amqp.Dial(rabbitMQURL)
		if err == nil {
			return conn, nil
		}
		time.Sleep(5 * time.Second)
	}
	return nil, err
}

type RabbitMQQueue struct {
	ch *amqp.Channel
}

func NewQueue(conn *amqp.Connection) (*RabbitMQQueue, error) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ channel: %v", err) // Exit here if channel creation fails
		return nil, err
	}

	return &RabbitMQQueue{ch}, err
}

func (q *RabbitMQQueue) DeclareQueue(queueName string) error {
	_, err := q.ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	return err
}

func (q *RabbitMQQueue) Publish(message string) error {
	return q.ch.Publish(
		"",                   // exchange
		"verification_queue", // routing key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

func (q *RabbitMQQueue) Consume(queueName string) (<-chan amqp.Delivery, error) {
	return q.ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
}
