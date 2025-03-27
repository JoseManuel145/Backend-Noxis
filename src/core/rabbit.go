package core

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ConnRabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Err        string
}

func GetRabbitMQ() *ConnRabbitMQ {
	error := ""
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error al cargar el archivo .env: %v", err)
	}

	rabbitMQURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASS"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)

	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		error = fmt.Sprintf("Error al conectar con RabbitMQ: %v", err)
		return &ConnRabbitMQ{Connection: nil, Channel: nil, Err: error}
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		error = fmt.Sprintf("Error al abrir un canal en RabbitMQ: %v", err)
		return &ConnRabbitMQ{Connection: nil, Channel: nil, Err: error}
	}

	return &ConnRabbitMQ{Connection: conn, Channel: ch, Err: error}
}
