package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"utils/requests"

	"data"
	"data/algorithmsData"
	"data/server"
	"net/url"
	"utils"
)

// Call the function to check errors, if any, print them
func checkAndPrintErrors(config data.Config) bool {
	errors := checkConfigValidity(config)
	if  errors != nil {
		for _, err := range errors {
			utils.Print(data.Red, err)
		}

		return true
	}

	return false
}

// Checks for configuration validity: a valid algorithm and Server field not empty
func checkConfigValidity(config data.Config) []string {
	var errors []string

	if _, exists := algorithmsData.LBAlgorithms[config.Algorithm]; !exists {
		errors = append(errors, "[+] Unable to find algorithm " + config.Algorithm)
	}

	if len(config.Servers) == 0 {
		errors = append(errors, "[+] No server found ")
	}

	return errors
}

// create an array of servers based on URLs in config files
func initializeServers(config data.Config) []server.Server{
	var servers []server.Server

	for _, serverStr := range config.Servers {
		parsedURL, _ := url.Parse(serverStr)
		servers = append(servers, server.Server{URL: parsedURL})
	}

	return servers
}

// if the configurations are valid, print them
func printJsonData(config data.Config) {
	utils.Print(data.Green, "[+] Using algorithm ", config.Algorithm)
	
	utils.Print(data.Green, "[+] Servers: ")
	for _, server := range config.Servers {
		fmt.Printf("\t- %s\n", server)
	}
		
	utils.Print(data.Green, "[+] HealthCheck: ", config.HealthCheck)
	utils.Print(data.Green, "[+] Logs: ", config.Logs)
}

// read a json file
func readJson(config *data.Config) {
	jsonFile, err := os.Open("./configs/config.json")
    
    if err != nil {
        fmt.Println(err)
	}

    defer jsonFile.Close()
    byteValue, _ := io.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &config)
}

func main() {
	var config data.Config
	var servers []server.Server

	algorithmsData.Init()

	utils.Print(data.Green, "===== Starting SilkRoute =====\n")
	utils.Print(data.Green, "[+] Reading config files")
	readJson(&config)

	if checkAndPrintErrors(config) { return }

	printJsonData(config)

	servers = initializeServers(config)
	requests.StartServer(config, servers)

}