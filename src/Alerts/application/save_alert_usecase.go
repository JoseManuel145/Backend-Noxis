package application

import (
	"Backend/src/Alerts/domain"
	"Backend/src/Alerts/domain/repositories"
	"fmt"
)

type SaveAlert struct {
	repo repositories.IAlertRepository
}

func NewSaveAlert(repo repositories.IAlertRepository) *SaveAlert {
	return &SaveAlert{repo: repo}
}

func (uc *SaveAlert) Execute(alert *domain.Alert) error {
	// Logica para elegir si guardar o no la alerta

	if alert.Sensor == "" || len(alert.Data) == 0 {
		return fmt.Errorf("alerta inválida: sensor y datos requeridos")
	}
	return uc.repo.SaveAlert(alert)
}
