package requests

import (
	"crypto/tls"
	"log"
	"net"
	"sync"
	"time"
	
	"data"
	"data/algorithmsData"
	"internals"
	"internals/healthCheck"
	"utils"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
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

// obtains the configuration for the tls certificate
func obtainCertificate(domain string, cache string) *tls.Config {
	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain),
		Cache:      autocert.DirCache(cache),
	}

	cfg := &tls.Config{
		GetCertificate: m.GetCertificate,
		NextProtos: []string{
			"http/1.1", acme.ALPNProto,
		},
	}

	return cfg
}

// handles request, sends a well formed request to a server
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

	err := fasthttp.DoTimeout(&ctx.Request, &ctx.Response, time.Second * 10)
	if err != nil {
		log.Print(err)
	}

	if config.Logs != "" {
		str := " " + ctx.RemoteIP().String() + " " + string(ctx.Method()) + " " + string(ctx.Path())
		log.Print(str)
		go utils.WriteLogs(str, config.Logs)
	}
}

// starts the listener on the port specified in host
func serveWithHttp(host string) {
	s := &fasthttp.Server{
		Handler:     requestHandler,
		ReadTimeout: time.Second * 2,
	}

	if err := s.ListenAndServe(host); err != nil {
		log.Fatal(err)
	}
}

// requests for a certificate and starts listener on 443 port
func serveWithHttps(host string) {
	var cfg *tls.Config = obtainCertificate(host, "./certs")

	ln, err := net.Listen("tcp4", host + ":443")
	if err != nil {
		log.Fatal(err)
	}

	var lnTls net.Listener = tls.NewListener(ln, cfg)

	if err := fasthttp.Serve(lnTls, requestHandler); err != nil {
		log.Fatal(err)
	}
}

// starts the listener to receive requests
func StartListener() {
	var config = data.GetConfig()

	if t := config.HealthCheck; t > 0 {
		go healthcheck.StartHealthCheckTimer(config.Servers, t, config.Dashboard < 0)
	}

	if len(config.Host) > 10 && config.Host[:10] == "localhost:" {
		serveWithHttp(config.Host)
	} else {
		serveWithHttps(config.Host)
	}
}