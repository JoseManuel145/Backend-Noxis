package core

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func StartConsumer() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error cargando .env: %v", err)
	}

	rabbitMQURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASS"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)

	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Fallo al conectar a RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Fallo al abrir un canal: %v", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"Sensores",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error al declarar el exchange: %v", err)
	}

	sensorBindings := map[string]string{
		"sensores/mq135/#":    "Calidad Aire MQ-135",
		"sensores/cjmcu811/#": "Carbono CJMCU-811",
		"sensores/mq7/#":      "Carbono MQ-7",
		"sensores/ky026/#":    "Flama KY-026",
		"sensores/mq5/#":      "Gas Natural MQ-5",
		"sensores/mq136/#":    "Hidrógeno MQ-136",
		"sensores/mq4/#":      "Metano MQ-4",
		"sensores/bme680/#":   "Presión-Temperatura-Humedad BME-680",
	}

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error al declarar la cola: %v", err)
	}

	for routingKey, sensorName := range sensorBindings {
		err = ch.QueueBind(
			q.Name,
			routingKey,
			"Sensores",
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("Error al enlazar la cola %s con el exchange Sensores: %v", sensorName, err)
		}
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("No se pudo consumir mensajes: %v", err)
	}

	fmt.Println("Consumidor iniciado y escuchando mensajes en el exchange 'Sensores'...")
	for msg := range msgs {
		fmt.Printf("[Mensaje Recibido] Routing Key: %s | Datos: %s\n", msg.RoutingKey, msg.Body)
	}
}
