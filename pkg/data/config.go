package data

import (
	"sync"

	"data/server"
)

// structure to hold json data
type ConfigRaw struct {
	Algorithm   string   `json:"algorithm"`
	Host        string   `json:"host"`
	Dashboard   int      `json:"dashboard"`
	Servers     []string `json:"servers"`
	HealthCheck int      `json:"healthCheck"`
	Logs        string   `json:"logs"`
	Prohibited  []string `json:"prohibited"`
}

// converts configurations from a raw format to the right format
func (rawConfig ConfigRaw) Cast() Config {
	var servers *server.Servers = GetConfig().Servers
	servers.Data = []server.ServerData{}

	for _, serverString := range rawConfig.Servers {
		var serverData = server.ServerData{URL: serverString, IsAlive: false}
		servers.Data = append(servers.Data, serverData)
	}

	return Config {
		Algorithm:   rawConfig.Algorithm,
		Host:        rawConfig.Host,
		Dashboard:   rawConfig.Dashboard,
		Servers:     servers,
		HealthCheck: rawConfig.HealthCheck,
		Logs:        rawConfig.Logs,
		Prohibited:  rawConfig.Prohibited,
	}
}

// singleton for the configuration

type Config struct {
	Path        string
	Algorithm   string
	Host        string
	Dashboard   int
	Servers     *server.Servers
	HealthCheck int
	Logs        string
	Prohibited  []string
}

var (
	configInstance *Config
	configOnce     sync.Once
)

func GetConfig() *Config {
	configOnce.Do(func() {
		configInstance = &Config{
			Servers: &server.Servers{
				Data: []server.ServerData{},
			},
		}
	})

	return configInstance
}

// Checks for configuration validity: a valid algorithm, Servers and Host field not empty
func (config Config) CheckValidity() []string {
	var errors []string

	switch config.Algorithm {
	case "rnd", "rr", "iph":
		break
	default:
		errors = append(errors, "[-] Set a valid algorithm (rnd, rr, iph)")
	}

	if len(config.Host) == 0 {
		errors = append(errors, "[-] Set a valid host ")
	}

	if len(config.Servers.Data) == 0 {
		errors = append(errors, "[-] No server found ")
	}

	return errors
}
