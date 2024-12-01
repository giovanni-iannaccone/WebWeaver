package webui

import (
	"fmt"
	"html/template"

	"data"
	"utils"

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
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}

	var addr string = ":" + fmt.Sprint(config.Dashboard)
	go fasthttp.ListenAndServe(addr, html)
}