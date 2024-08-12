package internals

import (
	"data"
	"data/algorithmsData"
)

// Checks for configuration validity: a valid algorithm and Server field not empty
func CheckConfigValidity(config data.Config) []string {
	var errors []string

	if _, exists := algorithmsData.LBAlgorithms[config.Algorithm]; !exists {
		errors = append(errors, "[+] Unable to find algorithm " + config.Algorithm)
	}

	if len(config.Servers) == 0 {
		errors = append(errors, "[+] No server found ")
	}

	return errors
}