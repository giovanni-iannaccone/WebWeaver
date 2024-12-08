package webui

import (
	"encoding/json"
	"fmt"
	"html/template"
	"path/filepath"

	"data"
	"data/server"
	"utils"

	"github.com/valyala/fasthttp"
)

type pageData struct {
	Servers []server.Server
	Port 	int
}

var tpl *template.Template

// parses the template
func Init() {
	tpl = template.Must(template.ParseGlob("webui/templates/index.html"))
}

// reads the configuration file and update the configurations
func hotReload(config *data.Config) {
	*config = utils.ReadAndParseJson(config.Path)
}

// executes the template
func idx(ctx *fasthttp.RequestCtx, pd pageData) {
	ctx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")

	err := tpl.ExecuteTemplate(ctx.Response.BodyWriter(), "index.html", pd)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

// renders the template on the given port
func RenderUI(config *data.Config) {
	var pd = pageData{config.Servers, config.Dashboard}

	html := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/":
			idx(ctx, pd)

		case "/hot-reload/":
			hotReload(config)
		
		case "/ws/":
			sendData(ctx, config.Servers)

		default:
			staticFileHandler(ctx)
		}
	}

	var addr string = ":" + fmt.Sprint(config.Dashboard)
	go fasthttp.ListenAndServe(addr, html)
}

// sends data to the webpage using a websocket
func sendData(ctx *fasthttp.RequestCtx, servers []server.Server) {
	var ws = data.GetWebSocket()
	ws.UpgradeToWS(ctx)

	data, err := json.Marshal(servers)
	if err != nil {
		ws.Conn.WriteMessage(1, data)
	}
}

// servers static files like css and js
func staticFileHandler(ctx *fasthttp.RequestCtx) {
	var file = string(ctx.Path())
	var staticDir = "webui/templates"
	var fullPath = filepath.Join(staticDir, file)

	ext := filepath.Ext(file)
	switch ext {
	case ".css":
		ctx.Response.Header.Set("Content-Type", "text/css")
	case ".js":
		ctx.Response.Header.Set("Content-Type", "application/javascript")
	default:
		ctx.Error("not found", fasthttp.StatusNotFound)
		return
	}

	fasthttp.ServeFile(ctx, fullPath)
}
