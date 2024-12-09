package healthcheck

import (
	"net"
	"time"

	"data"
	"data/server"
	"utils"
)

// checks if all servers are alive and notify observers
func HealthCheck(servers *server.Servers) {
	var updatedAnyServer bool = false

	for i := range servers.Data {
		previousState := servers.Data[i].IsAlive
		currentState := isServerAlive(servers.Data[i].URL)

		if previousState != currentState {
			servers.Data[i].IsAlive = currentState
			updatedAnyServer = true
		}
	}

	if updatedAnyServer {
		servers.NotifyObservers()
	}
}

// sends a request to the server to check if it is alive
func isServerAlive(url string) bool {
	var timeout time.Duration = 2 * time.Second

	conn, err := net.DialTimeout("tcp", url, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// prints servers status
func PrintHealthCheckStatus(servers *server.Servers) {
	for  i := range servers.Data {
		if servers.Data[i].IsAlive {
			utils.Print(data.Green, "[+] %s\t\talive\n", servers.Data[i].URL)
		} else {
			utils.Print(data.Yellow, "[!] %s\t\tNOT alive\n", servers.Data[i].URL)
		}
	}
}

// starts the health check timer, call the healthcheck function every time the timer expires
func StartHealthCheckTimer(servers *server.Servers, seconds int, printHealthCheckResult bool) {
	t := time.NewTicker(time.Second * time.Duration(seconds))
	defer t.Stop()

	for range t.C {
        HealthCheck(servers)
		if printHealthCheckResult {
        	PrintHealthCheckStatus(servers)
		}
    }
}