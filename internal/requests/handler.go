package requests

import (
	"data"
	"data/algorithmsData"
	"data/server"
	"internals"
	"internals/healthCheck"
	"utils"

	"log"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	config  *data.Config
	servers server.ServersData
	mu      sync.Mutex
)

// obtain next server based on the algorithm
func getNextServer(ip string) {

	algorithmsData.LBAlgorithms[config.Algorithm](servers, ip)

	for !servers.List[servers.Using].IsAlive {
		algorithmsData.LBAlgorithms[config.Algorithm](servers, ip)
	}
}

// handle request, send a well formed request to a server
func requestHandler(ctx *fasthttp.RequestCtx) {
	mu.Lock()
	getNextServer(ctx.RemoteIP().String())
	mu.Unlock()

	if internals.IsProhibited(config.Prohibited, ctx.Path()) {
		ctx.Error("404 not found", fasthttp.StatusNotFound)
		return
	}

	ctx.Request.SetHost(servers.List[servers.Using].URL.String())
	err := fasthttp.DoTimeout(&ctx.Request, &ctx.Response, time.Second*10)
	if err != nil {
		log.Print(err)
	}

	if config.Logs != "" {
		str := " " + ctx.RemoteIP().String() + " " + string(ctx.Method()) + " " + string(ctx.Path())
		log.Print(str)
		go utils.WriteLogs(str, config.Logs)
	}
}

// start the listener to receive requests
func StartListener(configurations *data.Config, serversList server.ServersData) {
	mu.Lock()
	config = configurations
	servers = serversList
	mu.Unlock()

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
