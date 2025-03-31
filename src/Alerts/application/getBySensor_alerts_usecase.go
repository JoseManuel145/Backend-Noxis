package application

import (
	"Backend/src/Alerts/domain"
	"Backend/src/Alerts/domain/repositories"
)

type GetBySensorAlerts struct {
	repo repositories.IAlertRepository
}

func NewGetBySensorAlert(repo repositories.IAlertRepository) *GetBySensorAlerts {
	return &GetBySensorAlerts{repo: repo}
}

func (uc *GetBySensorAlerts) Execute(sensor string) ([]domain.Alert, error) {
	return uc.repo.GetBySensor(sensor)
}
