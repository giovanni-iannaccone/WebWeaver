package data

import "algorithms"

var Algorithms map[string]func([]string, bool, bool)

func Init() {
	Algorithms = map[string]func([]string, bool, bool){
		"lc":  algorithms.LeastConnections,
		"rnd": algorithms.Random,
		"rr":  algorithms.RoundRobin,
		"wrr": algorithms.WeightedRoundRobin,
	}
}
