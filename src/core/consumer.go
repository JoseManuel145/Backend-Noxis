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

	// Definir argumentos para las colas
	args := amqp.Table{
		"x-max-length":  100,         // Máximo de 100 mensajes en la cola
		"x-message-ttl": 600000,      // Mensajes expiran después de 10 minutos
		"x-overflow":    "drop-head", // Si la cola está llena, se eliminan los mensajes más antiguos
		"x-queue-type":  "classic",   // Tipo de cola clásica
	}
	// Mapeo de colas a sus respectivas routing keys
	sensorBindings := map[string]string{
		"Calidad Aire MQ-135":                 "sensores.mq135.#",
		"Carbono CJMCU-811":                   "sensores.cjmcu811.#",
		"Carbono MQ-7":                        "sensores.mq7.#",
		"Flama KY-026":                        "sensores.flama.#",
		"Gas Natural MQ-5":                    "sensores.mq5.#",
		"Hidrogeno MQ-136":                    "sensores.mq136.#",
		"Metano MQ-4":                         "sensores.mq4.#",
		"Presion-Temperatura-Humedad BME-680": "sensores.bme680.#",
	}

	// Canal de goroutines para mantener consumidores corriendo
	done := make(chan bool)

	for queueName, routingKey := range sensorBindings {
		go func(queueName, routingKey string) {
			// Declarar la cola
			q, err := ch.QueueDeclare(
				queueName, // Nombre exacto de la cola
				true,      // Durable
				false,     // No auto-delete
				false,     // No exclusiva
				false,     // No-wait
				args,      // Argumentos de la cola
			)
			if err != nil {
				log.Fatalf("Error al declarar la cola %s: %v", queueName, err)
			}

			// Enlazar la cola con el exchange amq.topic
			err = ch.QueueBind(
				q.Name,
				routingKey,
				"amq.topic", // Exchange predefinido en RabbitMQ
				false,
				nil,
			)
			if err != nil {
				log.Fatalf("Error al enlazar la cola %s con el exchange amq.topic: %v", queueName, err)
			}

			// Consumir mensajes de la cola
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
				log.Fatalf("No se pudo consumir mensajes en la cola %s: %v", queueName, err)
			}

			fmt.Printf("Consumidor iniciado para la cola '%s', esperando mensajes...\n", queueName)
			for msg := range msgs {
				fmt.Printf("[Mensaje Recibido] Cola: %s | Routing Key: %s | Datos: %s\n", queueName, msg.RoutingKey, msg.Body)
			}
		}(queueName, routingKey)
	}

	<-done // Mantener el programa en ejecución
}
