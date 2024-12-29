package healthcheck

import (
	"net"
	"time"

	"data"
	"data/server"
	"utils"
)

const ACTIVE = true
const INACTIVE = false

func changeServersArray(servers *server.Servers,
	updatedActive []int,
	updatedInactive []int) {

		for i := range updatedActive {
			servers.Inactive = append(servers.Inactive, servers.Active[i])
			servers.Active = append(servers.Active[:i], servers.Active[i + 1:]...)
		}

		for i := range updatedInactive {
			servers.Active = append(servers.Active, servers.Inactive[i])
			servers.Inactive = append(servers.Inactive[:i], servers.Inactive[i + 1:]...)
		}
}

// checks if a specific list of servers are still in their state
func checkServers(servers []string, state bool) []int {
	var updatedServers []int

	for i := range servers {
		currentState := isServerAlive(servers[i])

		if currentState != state {
			updatedServers = append(updatedServers, i)
		}
	}

	return updatedServers
}

// checks if all servers are alive and notify observers
func HealthCheck(servers *server.Servers) {
	var updatedActiveServers []int
	var updatedInactiveServers []int

	updatedActiveServers = checkServers(servers.Active, ACTIVE)
	updatedInactiveServers = checkServers(servers.Inactive, INACTIVE)

	if len(updatedActiveServers) > 0 || len(updatedInactiveServers) > 0 {
		changeServersArray(servers, updatedActiveServers, updatedInactiveServers)
		servers.NotifyObservers()
	}
}

// sends a request to the server to check if it is alive
func isServerAlive(url string) bool {
	var timeout time.Duration = 2 * time.Second

	conn, err := net.DialTimeout("tcp", url, timeout)
	if err != nil {
		return INACTIVE
	}
	defer conn.Close()

	return ACTIVE
}

// prints servers status
func PrintHealthCheckStatus(servers *server.Servers) {
	var currentTime time.Time = time.Now()
    var formattedTime string = currentTime.Format(time.ANSIC)
	utils.Print(data.Blue, "\n\nHealthcheck at %s\n", formattedTime)

	for i := range servers.Inactive {
		utils.Print(data.Yellow, "[!] %s\t\tNOT alive\n", servers.Inactive[i])
	}

	for i := range servers.Active {
		utils.Print(data.Green, "[+] %s\t\t alive\n", servers.Active[i])
	}

	utils.Print(data.Reset, "\n\n")
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
