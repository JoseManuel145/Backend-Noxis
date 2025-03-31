package application

import (
	"Backend/src/Alerts/domain"
	"Backend/src/Alerts/domain/repositories"
)

type GetAllAlerts struct {
	repo repositories.IAlertRepository
}

func NewGetAllAlerts(repo repositories.IAlertRepository) *GetAllAlerts {
	return &GetAllAlerts{repo: repo}
}

func (uc *GetAllAlerts) Execute() ([]domain.Alert, error) {
	return uc.repo.GetAlerts()
}
