package algorithmsData

import (
	"algorithms"
	"data/server"
)

// Load Balancing Algorithms
var LBAlgorithms map[string]func([]server.Server) *server.Server

func Init() {
	LBAlgorithms = map[string]func([]server.Server) *server.Server{
		"iph": algorithms.IpHash,
		"lc":  algorithms.LeastConnections,
		"rnd": algorithms.Random,
		"rr":  algorithms.RoundRobin,
		"wrr": algorithms.WeightedRoundRobin,
	}
}
