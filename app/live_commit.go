package app

import (
	"github.com/gorilla/websocket"
	"github.com/xxlixin1993/LCS/logging"
	"net/http"
)

// Init WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// Don't check origin
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Handle WebSocket live commit
func LiveCommit(w http.ResponseWriter, r *http.Request, roomId uint32) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.ErrorF("Upgrade error(%s)", err)
		return
	}

	room := GetRoom(roomId)

	client := &Client{
		room: room,
		conn: conn,
		msg:  make(chan []byte, 256),
	}

	client.room.register <- client

	// TODO configure goroutine
	go client.writeGoroutine()
	go client.readGoroutine()
}
