package websocket

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
)

var websocketConnectionsLock sync.Mutex
var websocketConnections []*websocket.Conn

func registerWebsocketConnection(c *websocket.Conn) {
	websocketConnectionsLock.Lock()
	defer websocketConnectionsLock.Unlock()
	websocketConnections = append(websocketConnections, c)
}

func unregisterWebsocketConnection(c *websocket.Conn) bool {
	websocketConnectionsLock.Lock()
	defer websocketConnectionsLock.Unlock()
	for _, websocketConnection := range websocketConnections {
		if websocketConnection == c {
			websocketConnections = websocketConnections[:len(websocketConnections)-1]
			return true
		}
	}
	return false
}
