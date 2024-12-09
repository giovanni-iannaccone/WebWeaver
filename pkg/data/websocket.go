package data

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// singleton for the websocket
type WebSocket struct {
	Conn    *websocket.Conn
	upgrader websocket.Upgrader
}

var (
	instance *WebSocket
	once     sync.Once
)

func GetWebSocket() *WebSocket {
	once.Do(func() {
		instance = &WebSocket{
			upgrader: websocket.Upgrader{},
		}
	})
	return instance
}

// upgrades http connection to websocket
func (ws *WebSocket) UpgradeToWS(w http.ResponseWriter, r *http.Request) error {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	ws.Conn = conn
	return nil
}
