package healthcheck

import (
	"net"
	"net/url"
	"time"

	"data"
	"data/server"
	"utils"
)

// checks if all servers are alive	
func HealthCheck(servers []server.Server) {
	for i := range servers {
		servers[i].IsAlive = isServerAlive(servers[i].URL)
	}
}

// sends a request to the server to check if it is alive
func isServerAlive(u *url.URL) bool {
	var timeout time.Duration = 2 * time.Second

	conn, err := net.DialTimeout("tcp", u.String(), timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// prints servers status
func PrintHealthCheckStatus(servers []server.Server) {
	for _, s := range servers {
		if s.IsAlive {
			utils.Print(data.Green, "[+] %s\t\talive\n", s.URL.String())
		} else {
			utils.Print(data.Yellow, "[!] %s\t\tNOT alive\n", s.URL.String())
		}
	}
}

// starts the health check timer, call the healthcheck function every time the timer expires
func StartHealthCheckTimer(servers *[]server.Server, seconds int) {
	t := time.NewTicker(time.Second * time.Duration(seconds))
	defer t.Stop()

	for range t.C {
        HealthCheck(*servers)
        PrintHealthCheckStatus(*servers)
    }
}