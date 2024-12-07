package webui

import (
	"fmt"
	"html/template"
	"path/filepath"

	"data"
	"utils"

	//"github.com/gorilla/websocket"
	"github.com/valyala/fasthttp"
)

var tpl *template.Template

// parses the template
func Init() {
	tpl = template.Must(template.ParseGlob("webui/templates/index.html"))
}

// reads the configuration file and update the configurations
func hotReload(path string) data.Config {
	return utils.ReadAndParseJson(path)
}

// executes the template
func idx(ctx *fasthttp.RequestCtx, config data.Config) {
	ctx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")

	err := tpl.ExecuteTemplate(ctx.Response.BodyWriter(), "index.html", config.Servers)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

// renders the template on the given port
func RenderUI(config *data.Config) {
	html := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/":
			idx(ctx, *config)

		case "/hot-reload/":
			*config = hotReload(config.Path)

		default:
			staticFileHandler(ctx)
		}
	}

	var addr string = ":" + fmt.Sprint(config.Dashboard)
	go fasthttp.ListenAndServe(addr, html)
}

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
