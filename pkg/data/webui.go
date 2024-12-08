package data

import (
    "errors"
	"sync"

    "github.com/valyala/fasthttp"
	"github.com/fasthttp/websocket"
)

// singleton for the websocket
type WebSocket struct {
    Conn    *websocket.Conn
    upgrader websocket.FastHTTPUpgrader
}

var instance *WebSocket
var once sync.Once

func GetWebSocket() *WebSocket {
    once.Do(func() {
        instance = &WebSocket{}
    })

    return instance
}

// upgrades http connection to websocket
func (ws WebSocket) UpgradeToWS(r *fasthttp.RequestCtx) error {
    err := ws.upgrader.Upgrade(r, func(conn *websocket.Conn) {
        ws.Conn = conn
    })

    if err != nil {
        return errors.New("failed upgrading the websocket")
    }

    return nil
}