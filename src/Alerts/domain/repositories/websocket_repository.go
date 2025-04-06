package repositories

import "Backend/src/Alerts/domain"

type IWebSocketRepository interface {
	SendMessage(alert *domain.Alert) error
}
