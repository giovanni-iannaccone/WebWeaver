package lbalgorithms

import (
	"errors"
)

type LoadBalancer interface {
	NextServer(servers *[]string, ip string) int
}

// Load Balancing Algorithms
func NewLoadBalancer(algorithm string) (LoadBalancer, error) {
	switch algorithm {
	case "iph":
		return IpHashAlgorithm{}, nil
	case "rnd":
		return RandomAlgorithm{}, nil
	case "rr":
		return &RoundRobinAlgorithm{}, nil
	default:
		return nil, errors.New("nknown algorithm")
	}
}