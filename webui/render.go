package webui

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"data"
	"data/server"
	"utils"

	"github.com/gorilla/websocket"
)

var tpl *template.Template

// parses the HTML templates
func Init() {
	tpl = template.Must(template.ParseGlob("webui/templates/index.html"))
}

// reads the configuration file and updates the configurations
func hotReload(config *data.Config) {
	*config = utils.ReadAndParseJson(config.Path)
}

// renders the template with server data
func idx(w http.ResponseWriter, servers []server.ServerData) {
	if err := tpl.ExecuteTemplate(w, "index.html", servers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// starts the HTTP server and handles routing
func RenderUI() {
	var obs = make(chan bool)
	var config = data.GetConfig()
	config.Servers.AddObserver(obs)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		idx(w, config.Servers.Data)
	})
	http.HandleFunc("/hot-reload/", func(w http.ResponseWriter, r *http.Request) {
		hotReload(config)
	})
	http.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		if err := sendData(w, r, config.Servers, obs); err != nil {
			utils.Print(data.Red, "%s", err.Error())
		}
	})
	http.HandleFunc("/static/", staticFileHandler)

	addr := fmt.Sprintf(":%d", config.Dashboard)
	http.ListenAndServe(addr, nil)
}

// establishes a WebSocket connection and sends data to the client
func sendData(w http.ResponseWriter, r *http.Request, servers *server.Servers, obs chan bool) error {
	ws := data.GetWebSocket()

	err := ws.UpgradeToWS(w, r)
	if err != nil {
		return err
	}

	defer ws.Conn.Close()

	for range obs {
		bytes, _ := json.Marshal(servers)

		err := ws.Conn.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			return err
		}
	}

	return nil
}

// serves static files like CSS and JavaScript
func staticFileHandler(w http.ResponseWriter, r *http.Request) {
	var staticDir = "webui"
	var file = r.URL.Path
	var fullPath = filepath.Join(staticDir, filepath.Clean(file))

	var ext = filepath.Ext(file)
	switch ext {
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	default:
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, fullPath)
}
