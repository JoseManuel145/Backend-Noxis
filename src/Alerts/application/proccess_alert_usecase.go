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
			if EsPeligroso(alert.Sensor, alert.Data) {
				err := ps.saveRepo.repo.SaveAlert(&alert)
				if err != nil {
					log.Println("Error guardando alerta:", err)
				} else {
					log.Printf("Alerta peligrosa guardada: Sensor %s, Datos: %v", alert.Sensor, alert.Data)
				}
			} else {
				log.Printf("Alerta ignorada: Sensor %s, Datos seguros: %v", alert.Sensor, alert.Data)
			}
		}
	}
}
