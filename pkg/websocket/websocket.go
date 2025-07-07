package websocket

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
	wsInstance *WebSocket
	wsOnce     sync.Once
)

func GetWebSocket() *WebSocket {
	wsOnce.Do(func() {
		wsInstance = &WebSocket{
			upgrader: websocket.Upgrader{},
		}
	})
	return wsInstance
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
