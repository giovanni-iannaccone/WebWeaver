package lbalgorithms

import (
	"crypto/md5"
	"encoding/binary"
	"math/rand"
)

// returns the index of the next server based on requester ip hash
type IpHashAlgorithm struct{}

func (IpHashAlgorithm) NextServer(servers *[]string, ip string) int {
	var length int = len(*servers)
	var index int = 0

	if length > 0 { 
		hash := md5.Sum([]byte(ip))
		hashInt := binary.BigEndian.Uint32(hash[:4])
		index = int(hashInt % uint32(len(*servers)))
	} else {
		index = -1
	}

	return index
}

// returns the index of next server in the list
type RoundRobinAlgorithm struct {
	index int
}

func (a *RoundRobinAlgorithm) NextServer(servers *[]string, _ string) int {
	var length int = len(*servers)

	if length > 0 {
		a.index = (a.index + 1) % len(*servers)
	} else {
		a.index = -1
	}

	return a.index
}

// returns a random number
type RandomAlgorithm struct{}

func (RandomAlgorithm) NextServer(servers *[]string, _ string) int {
	var length int = len(*servers)
	var index int = 0

	if length > 0 {
		index = rand.Int() % len(*servers)
	} else {
		index = -1
	}

	return index
}
