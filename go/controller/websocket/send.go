package websocket

import (
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
	"github.com/mjarkk/whatsapp-dev/go/models"
)

func SendBytes(payload []byte) {
	websocketConnectionsLock.Lock()
	defer websocketConnectionsLock.Unlock()

	for _, conn := range websocketConnections {
		conn.WriteMessage(websocket.TextMessage, payload)
	}
}

func SendString(payload string) {
	SendBytes([]byte(payload))
}

func SendJSON(data interface{}) {
	payload, err := json.Marshal(data)
	if err == nil {
		SendBytes(payload)
	}
}

func SendMessage(message models.Message) {
	SendJSON(struct {
		Type    string         `json:"type"`
		Message models.Message `json:"message"`
	}{
		Type:    "message",
		Message: message,
	})
}
