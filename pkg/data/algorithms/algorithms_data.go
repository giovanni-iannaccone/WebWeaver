package algorithmsData

import (
	"errors"
	
	"data/server"
	algorithms "internals"
)

type LoadBalancer interface {
	NextServer(servers []server.Server, ip string) int
}

// Load Balancing Algorithms
func NewLoadBalancer(algorithm string) (LoadBalancer, error) {
	switch algorithm {
	case "iph":
		return algorithms.IpHashAlgorithm{}, nil
	case "rnd":
		return algorithms.RandomAlgorithm{}, nil
	case "rr":
		return &algorithms.RoundRobinAlgorithm{}, nil
	default:
		return nil, errors.New("algoritmo sconosciuto")
	}
}