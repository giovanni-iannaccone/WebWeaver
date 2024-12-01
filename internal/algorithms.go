package internals

import (
	"crypto/md5"
	"encoding/binary"
	"math/rand"
	"data/server"
)

// returns the index of the next server based on requester ip hash
type IpHashAlgorithm struct{}

func (IpHashAlgorithm) NextServer(servers []server.Server, ip string) int {
	hash := md5.Sum([]byte(ip))
	hashInt := binary.BigEndian.Uint32(hash[:4])

	return int(hashInt % uint32(len(servers)))
}

// returns the index of next server in the list
type RoundRobinAlgorithm struct {
	index int
}

func (a *RoundRobinAlgorithm) NextServer(servers []server.Server, _ string) int {
	a.index = (a.index + 1) % len(servers)
	return a.index
}

// returns a random number
type RandomAlgorithm struct{}

func (RandomAlgorithm) NextServer(servers []server.Server, _ string) int {
	return rand.Int() % len(servers)
}
