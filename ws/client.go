package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID   int
	Conn *websocket.Conn
}

var clients = make(map[int]*Client)

func GetClientsByUserID(userID int) []*Client {
	var subscribedClients []*Client

	for _, client := range clients {
		if client.ID == userID {
			subscribedClients = append(subscribedClients, client)
		}
	}

	return subscribedClients
}

func HandleWebSocket(c *gin.Context) {
	// Получите идентификатор пользователя из запроса или аутентификации
	userID := 123 // Замените на ваш код получения идентификатора пользователя

	// Установите соединение WebSocket
	conn, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
	if err != nil {
		// Обработайте ошибку установки соединения WebSocket
		return
	}

	// Создайте новый клиент WebSocket
	client := &Client{
		ID:   userID,
		Conn: conn,
	}

	// Добавьте клиента в глобальную карту
	clients[userID] = client

	// Обработайте сообщения от клиента, если это необходимо
	// Например, клиент может отправлять команды или запросы обновлений.

	// Удалите клиента из глобальной карты при закрытии соединения
	defer func() {
		delete(clients, userID)
	}()

	// Ваш код обработки сообщений от клиента
}
