package websocket

import (
	"log"

	"github.com/gofiber/contrib/websocket"
)

var EventsRoute = websocket.New(func(c *websocket.Conn) {
	registerWebsocketConnection(c)

	var err error
	for {
		_, _, err = c.ReadMessage()
		if err != nil {
			break
		}
	}

	if !unregisterWebsocketConnection(c) {
		log.Println("Was unable to close websocket connection, this should not happen and is a bug")
	}
})
