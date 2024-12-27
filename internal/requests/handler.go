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
	mu 		sync.Mutex
	using 	int
)

// obtains next server based on the algorithm
func getNextServer(ip string) {
	var config = data.GetConfig()

	lb, err := algorithmsData.NewLoadBalancer(config.Algorithm)
	if err != nil {
		utils.Print(data.Red, "%s", err.Error())
		return
	}

	for !config.Servers.Data[using].IsAlive {
		using = lb.NextServer(&config.Servers.Data, ip)
	}
}

// handle request, send a well formed request to a server
func requestHandler(ctx *fasthttp.RequestCtx) {
	var config = data.GetConfig()

	mu.Lock()
	getNextServer(ctx.RemoteIP().String())

	if internals.IsProhibited(config.Prohibited, ctx.Path()) {
		ctx.Error("404 not found", fasthttp.StatusNotFound)
		mu.Unlock()
		return
	}

	ctx.Request.SetHost(config.Servers.Data[using].URL)
	mu.Unlock()

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
func StartListener() {
	var config = data.GetConfig()

	if t := config.HealthCheck; t > 0 {
		go healthcheck.StartHealthCheckTimer(config.Servers, t, config.Dashboard < 0)
	}

	s := &fasthttp.Server{
		Handler:     requestHandler,
		ReadTimeout: time.Second * 2,
	}

	if err := s.ListenAndServe(config.Host); err != nil {
		log.Fatal(err)
	}
}