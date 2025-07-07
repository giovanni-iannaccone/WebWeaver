package requests

import (
	"crypto/tls"
	"log"
	"logs"
	"net"
	"security"
	"sync"
	"time"

	"config"
	"console"
	"lbalgorithms"

	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

var (
	currentServerIdx  int
	mutex             sync.Mutex
)

// creates and returns a TLS configuration with a certificate manager
func createTLSConfig(domain string, cacheDir string) *tls.Config {
	var manager = &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain),
		Cache:      autocert.DirCache(cacheDir),
	}

	return &tls.Config{
		GetCertificate: manager.GetCertificate,
		NextProtos:     []string{"http/1.1", acme.ALPNProto},
	}
}

// main request handler
func handleRequest(ctx *fasthttp.RequestCtx) {
	var config = config.GetConfig()

	if security.IsProhibited(config.Prohibited, ctx.Path()) {
		ctx.Error("404 not found", fasthttp.StatusNotFound)
		return
	}

	mutex.Lock()
	selectNextServer(ctx.RemoteIP().String())

	if currentServerIdx < 0 {
		ctx.Error("500 internal server error", fasthttp.StatusInternalServerError)
		mutex.Unlock()
		return 
	} else {
		redirectToServer(ctx, config.Servers.Active[currentServerIdx])			
	}

	mutex.Unlock()

	var err = fasthttp.DoTimeout(&ctx.Request, &ctx.Response, 10*time.Second)
	if err != nil {
		logRequest(config.Logs, nil, err)
	}

	if config.Logs != "" {
		logRequest(config.Logs, ctx, nil)
	}
}

// logs a request or an error message
func logRequest(logFile string, ctx *fasthttp.RequestCtx, err error) {
	var logMessage = ""

	if ctx != nil {
		logMessage = " " + ctx.RemoteIP().String() + " " + string(ctx.Method()) + " " + string(ctx.Path())
	} else if err != nil {
		logMessage = err.Error()
	}

	logs.WriteLogs(logMessage, logFile)
	log.Print(logMessage)
}

// handles the redirection to the selected server
func redirectToServer(ctx *fasthttp.RequestCtx, server string) {
	ctx.Request.SetHost(server)
}

// retrieves the next server based on the load balancing algorithm
func selectNextServer(clientIP string) {
	var config = config.GetConfig()

	loadBalancer, err := lbalgorithms.NewLoadBalancer(config.Algorithm)
	if err != nil {
		console.Print(console.Red, "%s", err.Error())
		return
	}

	currentServerIdx = loadBalancer.NextServer(&config.Servers.Active, clientIP)
}

// starts the HTTP server on the given host
func startHTTPServer(host string, handler func(ctx *fasthttp.RequestCtx)) {
	var server = &fasthttp.Server{
		Handler:     handler,
		ReadTimeout: 2 * time.Second,
	}

	if err := server.ListenAndServe(host); err != nil {
		log.Fatal(err)
	}
}

// starts the HTTPS server on the given host
func startHTTPSServer(host string, handler func(ctx *fasthttp.RequestCtx)) {
	var	tlsConfig *tls.Config = createTLSConfig(host, "./certs")

	var listener, err = net.Listen("tcp4", host + ":443")
	if err != nil {
		log.Fatal(err)
	}

	var tlsListener = tls.NewListener(listener, tlsConfig)

	if err := fasthttp.Serve(tlsListener, handler); err != nil {
		log.Fatal(err)
	}
}

// starts the appropriate listener based on the configuration
func StartListener() {
	var config = config.GetConfig()

	var handler func(ctx *fasthttp.RequestCtx) = handleRequest

	if  len(config.Host) > 10 && 
		(config.Host[:10] == "localhost:" ||
		config.Host[:10] == "127.0.0.1:") {

		startHTTPServer(config.Host, handler)
	} else {
		startHTTPSServer(config.Host, handler)
	}
}