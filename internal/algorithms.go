package internals

// Load Balancing Algorithms implementations

import (
	"data/server"
	"math/rand"
)

// Gives a server based on the hash of the ip asking
func IpHash(serverList server.ServersData) {

}

// Just gives the next server in the list
func RoundRobin(serverList server.ServersData) {
	serverList.Using = (serverList.Using + 1) % len(serverList.List)
}

// Gives a random server
func Random(serverList server.ServersData) {
	serverList.Using = rand.Int() % len(serverList.List)
}
