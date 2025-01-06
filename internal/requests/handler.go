package requests

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"data"
	"data/algorithmsData"
	"internals"
	"utils"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"github.com/valyala/fasthttp"
)

const SESSION_ID_COOKIE = "sessionId"

var (
	currentServerIdx  int
	mutex             sync.Mutex
)

var sessionMapping = make(map[string]string)

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

// generates a random session cookie value
func generateSessionCookie() string {
	const cookieLength = 32

	var randomBytes = make([]byte, cookieLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatal("Error generating random bytes: " + err.Error())
	}

	return base64.URLEncoding.EncodeToString(randomBytes)
}

// retrieves the server assigned to the current session
func getStickySessionServer(ctx *fasthttp.RequestCtx) string {
	var server, exists = sessionMapping[
		string(ctx.Request.Header.Cookie(SESSION_ID_COOKIE)),
	]

	if exists {
		return server
	}

	return ""
}

// redirects traffic based on previous associations or generates a cookie
func handleStickySessions(ctx *fasthttp.RequestCtx) {
	var config = data.GetConfig()

	if internals.IsProhibited(config.Prohibited, ctx.Path()) {
		ctx.Error("404 not found", fasthttp.StatusNotFound)
		return
	}

	if server := getStickySessionServer(ctx); server != "" {
		redirectToServer(ctx, server)

	} else {
		var cookie string = setSessionCookie(ctx)

		mutex.Lock()
		selectNextServer(ctx.RemoteIP().String())
		mutex.Unlock()

		if currentServerIdx < 0 {
			ctx.Error("500 internal server error", fasthttp.StatusInternalServerError)
		}

		var server = config.Servers.Active[currentServerIdx]

		sessionMapping[cookie] = server
		redirectToServer(ctx, server)
	}

	var err = fasthttp.DoTimeout(&ctx.Request, &ctx.Response, 10*time.Second)
	if err != nil {
		logRequest(config.Logs, nil, err)
	}

	if config.Logs != "" {
		logRequest(config.Logs, ctx, nil)
	}
}

// main request handler
func handleRequest(ctx *fasthttp.RequestCtx) {
	var config = data.GetConfig()

	if internals.IsProhibited(config.Prohibited, ctx.Path()) {
		ctx.Error("404 not found", fasthttp.StatusNotFound)
		return
	}

	mutex.Lock()
	redirectBasedOnAlgorithm(ctx)
	redirectToServer(ctx, config.Servers.Active[currentServerIdx])
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

	utils.WriteLogs(logMessage, logFile)
	log.Print(logMessage)
}

// handles the redirection to the selected server
func redirectToServer(ctx *fasthttp.RequestCtx, server string) {
	ctx.Request.SetHost(server)
}

// redirects traffic using the configured algorithm
func redirectBasedOnAlgorithm(ctx *fasthttp.RequestCtx) error {
	selectNextServer(ctx.RemoteIP().String())

	if currentServerIdx < 0 {
		ctx.Error("500 internal server error", fasthttp.StatusInternalServerError)
		return errors.New("no server alive")
	}

	return nil
}

// retrieves the next server based on the load balancing algorithm
func selectNextServer(clientIP string) {
	var config = data.GetConfig()

	loadBalancer, err := algorithmsData.NewLoadBalancer(config.Algorithm)
	if err != nil {
		utils.Print(data.Red, "%s", err.Error())
		return
	}

	currentServerIdx = loadBalancer.NextServer(&config.Servers.Active, clientIP)
}

// sets a session ID cookie for sticky sessions
func setSessionCookie(ctx *fasthttp.RequestCtx) string {
	var cookie = &fasthttp.Cookie{}
	var cookieValue string = generateSessionCookie()

	cookie.SetKey(SESSION_ID_COOKIE)
	cookie.SetValue(cookieValue)
	cookie.SetExpire(time.Now().Add(24 * time.Hour))
	cookie.SetSameSite(fasthttp.CookieSameSiteLaxMode)

	ctx.Response.Header.SetCookie(cookie)
	
	return cookieValue
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
	var config = data.GetConfig()

	var handler func(ctx *fasthttp.RequestCtx)

	if config.Sticky {
		handler = handleStickySessions
	} else {
		handler = handleRequest
	}

	if len(config.Host) > 10 && config.Host[:10] == "localhost:" {
		startHTTPServer(config.Host, handler)
	} else {
		startHTTPSServer(config.Host, handler)
	}
}