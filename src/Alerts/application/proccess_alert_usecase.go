package application

import (
	_ "Backend/src/Alerts/domain"
	"Backend/src/Alerts/domain/repositories"
	"log"
	"time"
)

type ProcessSensor struct {
	RabbitMQ repositories.IRabbitMQService
	saveRepo *SaveAlert
}

func NewProcessSensor(Rabbit repositories.IRabbitMQService, save *SaveAlert) *ProcessSensor {
	return &ProcessSensor{RabbitMQ: Rabbit, saveRepo: save}
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
			err := ps.saveRepo.Execute(&alert)
			if err != nil {
				log.Printf("Error procesando alerta: Sensor %s, Datos: %v, Error: %v", alert.Sensor, alert.Data, err)
			}
		}
	}
}
