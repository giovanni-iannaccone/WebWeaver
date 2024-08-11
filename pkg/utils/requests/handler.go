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

type ProxyHandler struct {
	Proxy *httputil.ReverseProxy
}

func (t *ProxyHandler) RoundTrip(request *http.Request) (*http.Response, error) {
	return http.DefaultTransport.RoundTrip(request)
}

func NewProxyHandler(destUrl *url.URL) *ProxyHandler {
	ph := ProxyHandler{
		Proxy: httputil.NewSingleHostReverseProxy(destUrl),
	}
	ph.Proxy.Transport = &ph
	return &ph
}

func (h *ProxyHandler) ProxyRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("> ProxyRequest, Client: %v, %v %v %v\n", r.RemoteAddr, r.Method, r.URL, r.Proto)

	h.Proxy.ServeHTTP(w, r)
}

func HandleRequest(config data.Config, servers []server.Server) {
	algorithmsData.Init()

	var svrAddr string = "localhost:8080"
	var svrBaseUrl string = "/"
	var nextServer *server.Server = algorithmsData.LBAlgorithms[config.Algorithm](servers)

	proxyHandler := NewProxyHandler(nextServer.URL)
	http.HandleFunc(svrBaseUrl, proxyHandler.ProxyRequest)

	err := http.ListenAndServe(svrAddr, nil)
	if err != nil {
		panic(err)
	}

}
