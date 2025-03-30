package application

import (
	"Backend/src/Alerts/domain"
	"Backend/src/Alerts/domain/repositories"
	"encoding/json"
)

type SendAlert struct {
	webSocketRepo repositories.IWebSocketRepository
}

func NewSendAlertUseCase(repo repositories.IWebSocketRepository) *SendAlert {
	return &SendAlert{webSocketRepo: repo}
}

func (uc *SendAlert) Execute(alert *domain.Alert) error {
	jsonMessage, err := json.Marshal(alert)
	if err != nil {
		return err
	}
	uc.webSocketRepo.SendMessage(jsonMessage)
	return nil
}
