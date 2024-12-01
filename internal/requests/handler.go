package requests

import (
	"log"
	"sync"
	"time"
	
	"data"
	"data/algorithmsData"
	"internals"
	"internals/healthCheck"
	"utils"

	"github.com/valyala/fasthttp"
)

var (
	config  *data.Config
	mu      sync.Mutex
	using 	int
)

// obtains next server based on the algorithm
func getNextServer(ip string) {
	lb, err := algorithmsData.NewLoadBalancer(config.Algorithm)
	if err != nil {
		utils.Print(data.Red, err.Error())
		return
	}

	for !config.Servers[using].IsAlive {
		using = lb.NextServer(config.Servers, ip)
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

	ctx.Request.SetHost(config.Servers[using].URL.String())
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
func StartListener(configurations *data.Config,) {
	mu.Lock()
	config = configurations
	mu.Unlock()

	if t := config.HealthCheck; t > 0 {
		go healthcheck.StartHealthCheckTimer(&config.Servers, t)
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
