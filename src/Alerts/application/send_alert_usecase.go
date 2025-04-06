package application

import (
	"Backend/src/Alerts/domain"
	"Backend/src/Alerts/domain/repositories"
	"log"
)

type SendAlert struct {
	webSocketRepo repositories.IWebSocketRepository
}

func NewSendAlertUseCase(repo repositories.IWebSocketRepository) *SendAlert {
	return &SendAlert{webSocketRepo: repo}
}

func (uc *SendAlert) Execute(alert *domain.Alert) error {
	err := uc.webSocketRepo.SendMessage(alert)
	if err != nil {
		log.Println("error en uc send_alert", err)
		return err
	}
	log.Printf("Sensor enviado por WEBSOCKET: Sensor %s, Datos: %v", alert.Sensor, alert.Data)
	return nil
}
