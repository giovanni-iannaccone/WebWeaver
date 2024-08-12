package main

import (
	"net/url"

	"data"
	"data/algorithmsData"
	"data/server"
	"internals"
	"internals/requests"
	"utils"
)

// Call the function to check errors, if any, print them
func checkAndPrintErrors(config data.Config) bool {
	errors := internals.CheckConfigValidity(config)
	if errors != nil {
		for _, err := range errors {
			utils.Print(data.Red, err)
		}

		return true
	}

	return false
}

// create an array of servers based on URLs in config files
func initializeServers(config data.Config) []server.Server {
	var servers []server.Server

	for _, serverStr := range config.Servers {
		parsedURL, _ := url.Parse(serverStr)
		servers = append(servers, server.Server{URL: parsedURL})
	}

	return servers
}

// if the configurations are valid, print them
func printJsonData(config data.Config) {
	utils.Print(data.Green, "[+] Using algorithm %s\n", config.Algorithm)

	utils.Print(data.Green, "[+] Servers: \n")
	for _, server := range config.Servers {
		utils.Print(data.Gray, "\t- %s\n", server)
	}

	utils.Print(data.Green, "[+] HealthCheck: %t\n", config.HealthCheck)
	utils.Print(data.Green, "[+] Logs: %s\n", config.Logs)

	utils.Print(data.Green, "[+] Prohibited: \n")
	for _, dir := range config.Prohibited {
		utils.Print(data.Gray, "\t- %s\n", dir)
	}
}

func main() {
	var config data.Config
	var servers []server.Server

	algorithmsData.Init()

	utils.Print(data.Green, "===== Starting WebWeaver =====\n")
	utils.Print(data.Green, "[+] Reading config files\n")

	err := utils.ReadJson(&config, "./configs/config.json")
	if err != nil {
		utils.Print(data.Red, err.Error())
		return 
	}

	if checkAndPrintErrors(config) {
		return
	}

	printJsonData(config)

	servers = initializeServers(config)
	requests.StartServer(config, servers)
}
