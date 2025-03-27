package adapters

import (
	"encoding/json"
	"fmt"
	"log"

	"Backend/src/Alerts/domain"
	"Backend/src/Alerts/domain/repositories"
	"Backend/src/core"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQAdapter maneja la conexión y consumo de colas de RabbitMQ
type RabbitMQAdapter struct {
	conn    *core.ConnRabbitMQ
	channel *amqp.Channel
}

// NewRabbitMQAdapter inicializa RabbitMQAdapter
func NewRabbitMQAdapter() repositories.IRabbitMQService {
	conn := core.GetRabbitMQ()
	if conn.Err != "" {
		log.Fatalf("Error al configurar RabbitMQ: %v", conn.Err)
	}
	return &RabbitMQAdapter{conn: conn, channel: conn.Channel}
}

// FetchReports obtiene los mensajes de RabbitMQ y los convierte en objetos Alert
func (r *RabbitMQAdapter) FetchReports() ([]domain.Alert, error) {
	// Definir los argumentos de la cola
	args := amqp.Table{
		"x-max-length":  100,         // Máximo de 100 mensajes
		"x-message-ttl": 600000,      // Mensajes expiran después de 10 minutos
		"x-overflow":    "drop-head", // Si está llena, borra los mensajes más viejos
		"x-queue-type":  "classic",   // Tipo de cola clásica
	}

	// Definir las colas y routing keys
	sensorBindings := map[string]string{
		"Calidad Aire MQ-135": "sensores.mq135.#",
		"Carbono CJMCU-811":   "sensores.cjmcu811.#",
		"Carbono MQ-7":        "sensores.mq7.#",
		"Flama KY-026":        "sensores.flama.#",
		"Gas Natural MQ-5":    "sensores.mq5.#",
		"Hidrogeno MQ-136":    "sensores.mq136.#",
		"Metano MQ-4":         "sensores.mq4.#",
		"BME-680":             "sensores.bme680.#",
	}

	var alerts []domain.Alert

	// Iniciar un consumidor por cada cola en goroutines
	for queueName, routingKey := range sensorBindings {
		go r.consumeQueue(queueName, routingKey, args, &alerts)
	}
	select {} // Mantiene el proceso en ejecución
}

// consumeQueue configura y consume mensajes de una cola específica
func (r *RabbitMQAdapter) consumeQueue(queueName, routingKey string, args amqp.Table, alerts *[]domain.Alert) {
	// Declarar la cola
	q, err := r.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		log.Fatalf("Error al declarar la cola %s: %v", queueName, err)
	}

	// Enlazar la cola con el exchange amq.topic
	err = r.channel.QueueBind(
		q.Name,
		routingKey,
		"amq.topic",
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error al enlazar la cola %s con el exchange amq.topic: %v", queueName, err)
	}

	// Consumir mensajes
	msgs, err := r.channel.Consume(
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
		// Decodificar el mensaje recibido en JSON y asignarlo a Data
		var data map[string]any
		if err := json.Unmarshal(msg.Body, &data); err != nil {
			log.Printf("Error al decodificar el mensaje de la cola %s: %v", queueName, err)
			continue
		}

		// Crear una nueva alerta con los datos decodificados
		alert := domain.Alert{
			Sensor: queueName, // Asignamos el nombre de la cola al campo Sensor
			Data:   data,      // Asignamos los datos del mensaje al campo Data
		}

		// Agregar la alerta al slice de alertas
		*alerts = append(*alerts, alert)

		// Imprimir los datos para depuración
		fmt.Printf("[Mensaje Recibido] Cola: %s | Routing Key: %s | Datos: %v\n", queueName, msg.RoutingKey, data)
	}
}
