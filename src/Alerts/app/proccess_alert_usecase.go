package application

import (
	_ "Backend/src/Alerts/domain"
	"Backend/src/Alerts/domain/repositories"
	"log"
	"time"
)

type ProcessSensor struct {
	RabbitMQ repositories.IRabbitMQService
}

func NewProcessSensor(Rabbit repositories.IRabbitMQService) *ProcessSensor {
	return &ProcessSensor{RabbitMQ: Rabbit}
}

func (ps *ProcessSensor) StartProcessingSensors() {
	for {
		time.Sleep(2 * time.Second)

		data, err := ps.RabbitMQ.FetchReports()
		if err != nil {
			log.Println("Error obteniendo datos del sensor:", err)
			continue
		}
		for _, alert := range data {
			log.Printf("%s", alert.Data)
		}
		//for _, domain.Alert := range sensorData {
		//	log.Printf("[Sensor] Tipo: %s, Valor: %.2f %s, Tiempo: %s\n")
		//}
	}
}
