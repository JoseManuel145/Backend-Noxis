package application

import (
	"Backend/src/Alerts/domain/repositories"
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
	log.Println("Iniciando el procesamiento de sensores...") // Verificar que la función se ejecuta

	for {
		log.Println("Esperando datos del sensor...") // Mensaje de espera
		time.Sleep(2 * time.Second)                  // Simular intervalos de procesamiento

		log.Println("Llamando a FetchReports...") // Depuración adicional
		data, err := ps.RabbitMQ.FetchReports()
		if err != nil {
			log.Println("Error obteniendo datos del sensor:", err)
			continue
		}
		if len(data) == 0 {
			log.Println("No se obtuvieron datos del sensor.")
			continue
		}

		log.Printf("Datos obtenidos en usecase: %v", data) // Depuración

		for _, alert := range data {
			log.Printf("Procesando alerta: %v", alert)
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
