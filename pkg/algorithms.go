package algorithms

// Load Balancing Algorithms implementations

import (
	"data/server"
	"math/rand"
)

func IpHash(servers []server.Server) *server.Server {
	return &servers[rand.Int()%len(servers)]
}

func LeastConnections(servers []server.Server) *server.Server {
	return &servers[rand.Int()%len(servers)]
}

func RoundRobin(servers []server.Server) *server.Server {
	return &servers[rand.Int()%len(servers)]
}

func Random(servers []server.Server) *server.Server {
	return &servers[rand.Int() % len(servers)]
}

func WeightedRoundRobin(servers []server.Server) *server.Server {
	return &servers[rand.Int()%len(servers)]
}
