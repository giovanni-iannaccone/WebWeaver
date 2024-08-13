package data

import (
	"data/algorithmsData"
)

// struct that holds configuration info
type Config struct {
	Algorithm   string   `json:"algorithm"`
	Host        string   `json:"host"`
	Servers     []string `json:"servers"`
	HealthCheck int      `json:"healthCheck"`
	Logs        string   `json:"logs"`
	Prohibited  []string `json:"prohibited"`
}

// Checks for configuration validity: a valid algorithm and Server field not empt
func (config Config) CheckValidity() []string {
	var errors []string

	if _, exists := algorithmsData.LBAlgorithms[config.Algorithm]; !exists {
		errors = append(errors, "[+] Unable to find algorithm "+config.Algorithm)
	}

	if len(config.Host) == 0 {
		errors = append(errors, "[+] Set a valid host ")
	}

	if len(config.Servers) == 0 {
		errors = append(errors, "[+] No server found ")
	}

	return errors
}
