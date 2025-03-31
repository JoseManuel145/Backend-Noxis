package application

import (
	"Backend/src/Alerts/domain/repositories"
	"fmt"
	"log"
	"time"
)

type ProcessSensor struct {
	RabbitMQ repositories.IRabbitMQService
	saveRepo *SaveAlert
	sendMsg  *SendAlert
}

func NewProcessSensor(Rabbit repositories.IRabbitMQService, save *SaveAlert, send *SendAlert) *ProcessSensor {
	return &ProcessSensor{RabbitMQ: Rabbit, saveRepo: save, sendMsg: send}
}

func (ps *ProcessSensor) StartProcessingSensors() {
	fmt.Println("Iniciando el procesamiento de sensores...") // Verificar que la función se ejecuta

	for {
		time.Sleep(2 * time.Second) // Simular intervalos de procesamiento

		data, err := ps.RabbitMQ.FetchReports()
		if err != nil {
			log.Println("Error obteniendo datos del sensor:", err)
			continue
		}
		if len(data) == 0 {
			continue
		}

		fmt.Println("Datos obtenidos en usecase:", data)

		for _, alert := range data {
			err := ps.saveRepo.Execute(&alert)
			if err != nil {
				log.Printf("Error procesando alerta: Sensor %s, Datos: %v, Error: %v", alert.Sensor, alert.Data, err)
				continue
			}
			log.Printf("Alerta procesada correctamente: Sensor %s, Datos: %v", alert.Sensor, alert.Data)

			err = ps.sendMsg.Execute(&alert)
			if err != nil {
				log.Printf("Error enviando alerta al websocket: Sensor %s, Datos: %v, Error: %v", alert.Sensor, alert.Data, err)
			}
		}
	}
}
