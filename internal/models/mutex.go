package models

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Storing the clients connected to the server in a map
var clients = make(map[string]*websocket.Conn)

// Mutex to be able to lock before writing to the map
var mu sync.Mutex

func GetMux() *sync.Mutex {
	return &mu
}

func GetClientMap() map[string]*websocket.Conn {
	return clients
}
