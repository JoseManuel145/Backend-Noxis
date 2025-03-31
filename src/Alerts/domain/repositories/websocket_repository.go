package repositories

type IWebSocketRepository interface {
	SendMessage(message []byte)
}
