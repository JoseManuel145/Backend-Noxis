package repositories

type IWebSocketRepository interface {
	SendMessage(sensor string, message []byte)
}
