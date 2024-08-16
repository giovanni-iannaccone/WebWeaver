package requests

import (
	"data"
	"data/algorithmsData"
	"data/server"
	"internals"
	"internals/healthCheck"
	"utils"

	"log"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	config data.Config
	servers server.ServersData
)

func getNextServer() {
	algorithmsData.LBAlgorithms[config.Algorithm](servers)

	for !servers.List[servers.Using].IsAlive {
		algorithmsData.LBAlgorithms[config.Algorithm](servers)
	}
}

// Handle requests, call the function to determine the server, redirect here the request and write logs
func requestHandler(ctx *fasthttp.RequestCtx) {
	var str string
	getNextServer()

	if internals.IsProhibited(config.Prohibited, ctx.Path()) {
		ctx.Error("404 not found", fasthttp.StatusNotFound)

	} else {
		ctx.Request.SetHost(
			servers.List[servers.Using].URL.String(),
		)

		err := fasthttp.DoTimeout(&ctx.Request, &ctx.Response, time.Second*10)
		if err != nil {
			log.Print(err)
		}
	}

	if config.Logs != "" {
		str = " " + ctx.RemoteIP().String() + " " + string(ctx.Method()) + " " + string(ctx.Path())
		log.Print(str)
		go utils.WriteLogs(str, config.Logs)
	}
}

// "Main" server function, define globals, register the handler and listen to incoming requests
func StartListener(configurations data.Config, serversList server.ServersData) {
	algorithmsData.Init()

	config = configurations
	servers = serversList

	if t := config.HealthCheck; t > 0 {
		go healthcheck.StartHealthCheckTimer(t, &serversList)
	}

	s := &fasthttp.Server{
		Handler:     requestHandler,
		ReadTimeout: time.Second * 2,
	}

	if err := s.ListenAndServe(config.Host); err != nil {
		log.Fatal(err)
	}

	select {}
}
