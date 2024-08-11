package requests

import (
	"data"
	"data/algorithmsData"
	"data/server"

	"log"
	"time"

	"github.com/valyala/fasthttp"
)

var config data.Config
var servers []server.Server

// Handle requests, call the function to determine the server and redirect here the request
func requestHandler(ctx *fasthttp.RequestCtx) {
	nextServer := algorithmsData.LBAlgorithms[config.Algorithm](servers)

	ctx.Request.SetHost(nextServer.URL.String())
	err := fasthttp.DoTimeout(&ctx.Request, &ctx.Response, time.Second*10)
	if err != nil {
		log.Print(err)
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
