package repositories

import (
	"Backend/src/Alerts/domain"
)

type IRabbitMQService interface {
	FetchReports() ([]domain.Alert, error)
}
