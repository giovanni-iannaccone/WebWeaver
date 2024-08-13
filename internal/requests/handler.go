package requests

import (
	"data"
	"data/algorithmsData"
	"data/server"
	"internals"
	"net"
	"net/url"
	"utils"

	"log"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	config data.Config
	servers server.ServersData
)

// Check if server is alive
func healthCheck() {
	for _, s := range servers.List {
		s.IsAlive = isServerAlive(s.URL)
	}
}

// Send a request to the server to check if it is alive
func isServerAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// Print server status
func printHealthCheckStatus() {
	for _, s := range servers.List {
		if s.IsAlive {
			utils.Print(data.Green, "[+] %s is alive", s.URL.String())
		} else {
			utils.Print(data.Red, "[!] %s is NOT alive", s.URL.String())
		}
	}
}

// Handle requests, call the function to determine the server, redirect here the request and write logs
func requestHandler(ctx *fasthttp.RequestCtx) {
	var str string
	algorithmsData.LBAlgorithms[config.Algorithm](servers)

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
		str = " " + ctx.Request.String()
		log.Print(str)
		go utils.WriteLogs(str, config.Logs)
	}
}

// Start the health check timer, call the healthcheck function every time the timer expires
func startHealthCheckTimer() {
	t := time.NewTicker(
		time.Second * time.Duration(config.HealthCheck),
	)
	defer t.Stop()

	done := make(chan bool)

	go func() {
		for {
			select {
			case <-t.C:
				healthCheck()
				printHealthCheckStatus()
			case <-done:
				return
			}
		}
	}()

	done <- true
}

// "Main" server function, define globals, register the handler and listen to incoming requests
func StartListener(configurations data.Config, serversList server.ServersData) {
	algorithmsData.Init()

	config = configurations
	servers = serversList

	s := &fasthttp.Server{
		Handler:     requestHandler,
		ReadTimeout: time.Second * 2,
	}

	if err := s.ListenAndServe(config.Host); err != nil {
		log.Fatal(err)
	}

	if config.HealthCheck > 0 {
		startHealthCheckTimer()
	}
}
