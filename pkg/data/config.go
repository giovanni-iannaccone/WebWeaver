package data

import (
	"net/url"

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
	var servers []server.Server

	for _, serverString := range rawConfig.Servers {
		parsedURL, _ := url.Parse(serverString)
		servers = append(servers, server.Server{URL: parsedURL, IsAlive: false})
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

// the final struct used by the program
type Config struct {
	Path        string
	Algorithm   string
	Host        string
	Dashboard   int
	Servers     []server.Server
	HealthCheck int
	Logs        string
	Prohibited  []string
}

// Checks for configuration validity: a valid algorithm and Server field not empt
func (config Config) CheckValidity() []string {
	var errors []string

	if len(config.Host) == 0 {
		errors = append(errors, "[-] Set a valid host ")
	}

	if len(config.Servers) == 0 {
		errors = append(errors, "[-] No server found ")
	}

	return errors
}
