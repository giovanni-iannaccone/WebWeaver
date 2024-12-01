package webui

import (
	"fmt"
	"html/template"

	"data"
	"utils"

	"github.com/valyala/fasthttp"
)

var tpl *template.Template

// Parses the template
func Init() {
    tpl = template.Must(template.ParseGlob("webui/templates/index.html"))
}

// Reads the configuration file and update the configurations
func hotReload(config *data.Config) {
	utils.ReadJson(config, config.Path)
}

// Executes the template
func idx(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")

	err := tpl.ExecuteTemplate(ctx.Response.BodyWriter(), "index.html", nil)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

// Renders the template on the given port
func RenderUI(config *data.Config) {
	html := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/":
			idx(ctx)
		
		case "/hot-reload/":
			hotReload(config)

		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}

	var addr string = ":" + fmt.Sprint(config.Dashboard)
	go fasthttp.ListenAndServe(addr, html)
}