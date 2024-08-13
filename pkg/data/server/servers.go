package server

import (
	"net/url"
)

type Server struct {
	URL     *url.URL
	IsAlive bool
}

type ServersData struct {
	List  []Server
	Using int
}
