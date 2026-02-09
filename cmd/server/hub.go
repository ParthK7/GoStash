package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Hub struct {
	data       <-chan string
	watcherMap map[*websocket.Conn]bool
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
}

func NewHub(logchan <-chan string) *Hub {
	return &Hub{
		data:       logchan,
		watcherMap: make(map[*websocket.Conn]bool),
		register:   make(chan *websocket.Conn, 10),
		unregister: make(chan *websocket.Conn, 10),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case connection, ok := <-h.register:
			if !ok {
				fmt.Printf("nil")
			}
			h.watcherMap[connection] = true

		case connection, ok := <-h.unregister:
			if !ok {
				fmt.Printf("nil")
			}
			delete(h.watcherMap, connection)

		case request, ok := <-h.data:
			if !ok {
				fmt.Printf("nil")
			}

			for connection := range h.watcherMap {
				err := connection.WriteMessage(websocket.TextMessage, []byte(request)) //websocket write function syntax.
				if err != nil {
					delete(h.watcherMap, connection)
					connection.Close()
				}
			}
		}
	}
}
