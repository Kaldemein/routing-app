package queue

import (
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectToRabbitMQ(rabbitMQURL string, retryCount int, retryInterval time.Duration) (*amqp.Connection, error) {
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
