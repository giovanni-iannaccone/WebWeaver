package algorithmsData

import (
	"data/server"
	algorithms "internals"
)

// Load Balancing Algorithms
var LBAlgorithms map[string]func(server.ServersData, string)

func Init() {
	LBAlgorithms = map[string]func(server.ServersData, string){
		"iph": algorithms.IpHash,
		"rnd": algorithms.Random,
		"rr":  algorithms.RoundRobin,
	}
}
