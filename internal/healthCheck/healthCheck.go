package healthcheck

import (
	"net"
	"net/url"
	"time"

	"data"
	"data/server"
	"utils"
)

// Check if server is alive	
func HealthCheck(servers *server.ServersData) {
	for i, s := range servers.List {
		servers.List[i].IsAlive = isServerAlive(s.URL)
	}
}

// Send a request to the server to check if it is alive
func isServerAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.String(), timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// PrintHealthCheckStatus print servers status
func PrintHealthCheckStatus(servers *server.ServersData) {
	for _, s := range servers.List {
		if s.IsAlive {
			utils.Print(data.Green, "[+] %s is alive\n", s.URL.String())
		} else {
			utils.Print(data.Red, "[!] %s is NOT alive\n", s.URL.String())
		}
	}
}

// Start the health check timer, call the healthcheck function every time the timer expires
func StartHealthCheckTimer(seconds int, servers *server.ServersData) {
	t := time.NewTicker(time.Second * time.Duration(seconds))
	defer t.Stop()

	for range t.C {
        HealthCheck(servers)
        PrintHealthCheckStatus(servers)
    }

}
