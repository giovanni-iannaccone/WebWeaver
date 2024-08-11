package server

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Server struct {
	URL          *url.URL
	IsAlive      bool
	Mutex        sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}
