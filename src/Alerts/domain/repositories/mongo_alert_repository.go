package repositories

import "Backend/src/Alerts/domain"

type IAlertRepository interface {
	SaveAlert(alert *domain.Alert) error
	GetBySensor(sensor string) ([]domain.Alert, error)
	GetAlerts() ([]domain.Alert, error)
}
