package internals

// Load Balancing Algorithms implementations

import (
	"crypto/md5"
	"encoding/binary"
	"math/rand"

	"data/server"
)

// Gives a server based on the hash of the ip asking
func IpHash(serverList server.ServersData, ipAddress string) {
	hash := md5.Sum([]byte(ipAddress))
	hashInt := binary.BigEndian.Uint32(hash[:4])

	serverList.Using = int(hashInt % uint32(len(serverList.List)))
}

// Just gives the next server in the list
func RoundRobin(serverList server.ServersData, _ string) {
	serverList.Using = (serverList.Using + 1) % len(serverList.List)
}

// Gives a random server
func Random(serverList server.ServersData, _ string) {
	serverList.Using = rand.Int() % len(serverList.List)
}
