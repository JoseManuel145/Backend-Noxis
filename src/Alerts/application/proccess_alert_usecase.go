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
			//implementar logica para decidir cuando guardar y cuando no
			ps.saveRepo.repo.SaveAlert(&alert)
			log.Printf("%s", alert.Data)
		}
	}
}
