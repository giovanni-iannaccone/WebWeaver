package requests

import (
	"data"
	"data/algorithmsData"
	"data/server"
	"internals"
	"utils"

	"log"
	"time"

	"github.com/valyala/fasthttp"
)

var config data.Config
var servers []server.Server

// Handle requests, call the function to determine the server, redirect here the request and write logs
func requestHandler(ctx *fasthttp.RequestCtx) {
	var str string
	var nextServer *server.Server = algorithmsData.LBAlgorithms[config.Algorithm](servers)

	if internals.IsProhibited(config.Prohibited, ctx.Path()) {
		ctx.Error("404 not found", fasthttp.StatusNotFound)

	} else {
		ctx.Request.SetHost(nextServer.URL.String())
		err := fasthttp.DoTimeout(&ctx.Request, &ctx.Response, time.Second*10)
		if err != nil {
			log.Print(err)
		}
	}

	if config.Logs != "" {
		str = " " + ctx.Request.String()
		log.Print(str)
		go utils.WriteLogs(str, config.Logs)
	}

}

// "Main" function, define globals, register the handler and listen to incoming requests
func StartServer(configurations data.Config, serversList []server.Server) {
	algorithmsData.Init()

	config = configurations
	servers = serversList

	s := &fasthttp.Server{
		Handler:     requestHandler,
		ReadTimeout: time.Second * 2,
	}

	if err := s.ListenAndServe("localhost:8080"); err != nil {
		log.Fatal(err)
	}
}
