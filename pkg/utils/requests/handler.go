package requests

import (
	"data"
	"data/algorithmsData"
	"data/server"

	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var config data.Config
var servers []server.Server

type ProxyHandler struct {
	Proxy *httputil.ReverseProxy
}

// New connections handler, get the next server based on algorithm and redirect the request
func Handler(w http.ResponseWriter, r *http.Request) {
	nextServer := algorithmsData.LBAlgorithms[config.Algorithm](servers)
	proxyHandler := NewProxyHandler(nextServer.URL)

	proxyHandler.ProxyRequest(w, r)
}

// Return a new proxy handler
func NewProxyHandler(destUrl *url.URL) *ProxyHandler {
	return &ProxyHandler{
		Proxy: httputil.NewSingleHostReverseProxy(destUrl),
	}
}

// Log and serve http request
func (h *ProxyHandler) ProxyRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("> ProxyRequest, Client: %v, %v %v %v\n", r.RemoteAddr, r.Method, r.URL, r.Proto)
	h.Proxy.ServeHTTP(w, r)
}

// "Main" function, define globals, register the handler and listen to incoming requests
func StartServer(configurations data.Config, serversList []server.Server) {
	algorithmsData.Init()

	config = configurations
	servers = serversList

	http.HandleFunc("/", Handler)

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		panic(err)
	}
}
