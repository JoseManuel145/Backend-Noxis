package adapters

import (
	"Backend/src/Alerts/domain"
	"Backend/src/Alerts/domain/repositories"
	"Backend/src/core"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

const batchSize = 5

type RabbitMQAdapter struct {
	conn         *core.ConnRabbitMQ
	channel      *amqp.Channel
	sensorQueues map[string]string            // Mapa de colas y sus routing keys
	consumers    map[string]chan domain.Alert // Mapa de consumidores activos por cola
}

// NewRabbitMQAdapter inicializa RabbitMQAdapter y declara colas una sola vez
func NewRabbitMQAdapter() repositories.IRabbitMQService {
	conn := core.GetRabbitMQ()
	if conn.Err != "" {
		log.Fatalf("Error al configurar RabbitMQ: %v", conn.Err)
	}

	adapter := &RabbitMQAdapter{
		conn:    conn,
		channel: conn.Channel,
		sensorQueues: map[string]string{
			"Calidad Aire MQ-135": "sensores.mq135.#",
			"Carbono CJMCU-811":   "sensores.cjmcu811.#",
			"Carbono MQ-7":        "sensores.mq7.#",
			"Flama KY-026":        "sensores.flama.#",
			"Gas Natural MQ-5":    "sensores.mq5.#",
			"Hidrogeno MQ-136":    "sensores.mq136.#",
			"Metano MQ-4":         "sensores.mq4.#",
			"BME-680":             "sensores.bme680.#",
		},
		consumers: make(map[string]chan domain.Alert),
	}

	adapter.setupQueues()    // Declara las colas solo una vez
	adapter.startConsumers() // Inicia los consumidores solo una vez
	return adapter
}

// setupQueues declara las colas solo una vez al inicio
func (r *RabbitMQAdapter) setupQueues() {
	args := amqp.Table{
		"x-max-length":  100,
		"x-message-ttl": 600000,
		"x-overflow":    "drop-head",
		"x-queue-type":  "classic",
	}

	for queueName, routingKey := range r.sensorQueues {
		q, err := r.channel.QueueDeclare(queueName, true, false, false, false, args)
		if err != nil {
			log.Fatalf("Error al declarar la cola %s: %v", queueName, err)
		}

		err = r.channel.QueueBind(q.Name, routingKey, "amq.topic", false, nil)
		if err != nil {
			log.Fatalf("Error al enlazar la cola %s con el exchange amq.topic: %v", queueName, err)
		}

		fmt.Printf("Cola '%s' configurada correctamente.\n", queueName)
	}
}

// startConsumers inicia los consumidores una vez por cada cola, sin volver a declararlas
func (r *RabbitMQAdapter) startConsumers() {
	for queueName := range r.sensorQueues {
		// Crear un canal por cada consumidor
		dataChan := make(chan domain.Alert, batchSize)
		r.consumers[queueName] = dataChan

		// Iniciar un goroutine para consumir mensajes en cada cola
		go r.consumeQueue(queueName, dataChan)
	}
}

// consumeQueue consume los mensajes de una cola específica sin declararla nuevamente
func (r *RabbitMQAdapter) consumeQueue(queueName string, dataChan chan domain.Alert) {
	msgs, err := r.channel.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("No se pudo consumir mensajes en la cola %s: %v", queueName, err)
	}

	for msg := range msgs {
		var data map[string]any
		if err := json.Unmarshal(msg.Body, &data); err != nil {
			log.Printf("Error al decodificar el mensaje de la cola %s: %v", queueName, err)
			continue
		}

		alert := domain.Alert{
			Sensor: queueName,
			Data:   data,
		}

		dataChan <- alert
		log.Printf("Alerta agregada: %v", alert)
	}
}

func (r *RabbitMQAdapter) FetchReports() ([]domain.Alert, error) {
	var alerts []domain.Alert
	var mu sync.Mutex
	dataChan := make(chan domain.Alert, batchSize)
	var wg sync.WaitGroup

	log.Println("Iniciando FetchReports...") // Depuración

	// Iniciar una goroutine por cada cola (ya declarada previamente)
	for queueName := range r.sensorQueues {
		wg.Add(1)
		go func(queue string) {
			defer wg.Done()
			// Consumir mensajes directamente desde el mapa de consumidores
			consumerChan, exists := r.consumers[queue]
			if !exists {
				log.Printf("No existe un consumidor para la cola %s", queue)
				return
			}
			for alert := range consumerChan {
				log.Printf("Alerta recibida en FetchReports: %v", alert) // Depuración
				dataChan <- alert
			}
		}(queueName)
	}

	// Esperar a que todas las goroutines terminen antes de cerrar el canal
	go func() {
		wg.Wait()
		log.Println("Cerrando canal de datos en FetchReports...") // Depuración
		close(dataChan)                                           // Cerrar el canal después de que todas las goroutines terminen
	}()

	// Acumular datos hasta completar batchSize
	for alert := range dataChan {
		mu.Lock()
		alerts = append(alerts, alert)
		mu.Unlock()

		if len(alerts) >= batchSize {
			log.Printf("Se alcanzó el tamaño del lote: %d", batchSize) // Depuración
			break
		}
	}

	// Si las alertas obtenidas son suficientes, retornar
	if len(alerts) == 0 {
		log.Println("No se obtuvieron alertas en FetchReports.") // Depuración
		return nil, fmt.Errorf("no se obtuvieron alertas")
	}
	log.Printf("Alertas obtenidas en FetchReports: %v", alerts) // Depuración
	return alerts, nil
}
