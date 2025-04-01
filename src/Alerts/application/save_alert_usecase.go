package application

import (
	"Backend/src/Alerts/domain"
	"Backend/src/Alerts/domain/repositories"
	"fmt"
	"log"
)

type SaveAlert struct {
	repo repositories.IAlertRepository
}

func NewSaveAlert(repo repositories.IAlertRepository) *SaveAlert {
	return &SaveAlert{repo: repo}
}

func (uc *SaveAlert) Execute(alert *domain.Alert) error {
	if !EsPeligroso(alert.Sensor, alert.Data) {
		log.Printf("Alerta ignorada: Sensor %s, Datos seguros: %v", alert.Sensor, alert.Data)
		return nil
	}

	err := uc.repo.SaveAlert(alert)
	if err != nil {
		log.Printf("Error guardando alerta: Sensor %s, Datos: %v, Error: %v", alert.Sensor, alert.Data, err)
		return fmt.Errorf("no se pudo guardar la alerta: %w", err)
	}

	log.Printf("Alerta peligrosa guardada: Sensor %s, Datos: %v", alert.Sensor, alert.Data)
	return nil
}
