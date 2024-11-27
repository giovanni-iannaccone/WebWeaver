package webui

import (
	"fmt"
	"html/template"

	"github.com/valyala/fasthttp"
)

// Executes the template
func idx(ctx *fasthttp.RequestCtx, tpl *template.Template) {
	ctx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")

	err := tpl.ExecuteTemplate(ctx.Response.BodyWriter(), "index.html", nil)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

// Renders the template on the given port
func RenderUI(port uint16) {
	tpl := template.Must(template.ParseFiles("webui/templates/index.html"))

	html := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/":
			idx(ctx, tpl)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}

	addr := ":" + fmt.Sprint(port)

	go fasthttp.ListenAndServe(addr, html)
}