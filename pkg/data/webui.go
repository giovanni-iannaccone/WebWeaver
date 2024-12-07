package data

import (
	"sync"

	"github.com/gorilla/websocket"
)

// singleton for the websocket
type WebSocket websocket.Conn

var instance *WebSocket
var once sync.Once

func GetConfigManager() *WebSocket {
    once.Do(func() {
        instance = &WebSocket{}
    })

    return instance
}