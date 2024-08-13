package main

import (
	"net/url"

	"data"
	"data/algorithmsData"
	"data/server"
	"internals/requests"
	"utils"
)

// Call the function to check errors, if any, print them
func checkAndPrintErrors(config data.Config) bool {
	errors := config.CheckValidity()
	if errors != nil {
		for _, err := range errors {
			utils.Print(data.Red, err)
		}

		return true
	}

	return false
}

// create an array of servers based on URLs in config files
func initializeServers(urlList []string, servers server.ServersData) {
	for _, serverStr := range urlList {
		parsedURL, _ := url.Parse(serverStr)
		servers.List = append(servers.List, server.Server{URL: parsedURL})
	}
}

// if the configurations are valid, print them
func printJsonData(config data.Config) {
	utils.Print(data.Green, "[+] Using algorithm %s\n", config.Algorithm)
	utils.Print(data.Green, "[+] Host: %s\n", config.Host)

	utils.Print(data.Green, "[+] Servers: \n")
	for _, server := range config.Servers {
		utils.Print(data.Gray, "\t- %s\n", server)
	}

	utils.Print(data.Green, "[+] HealthCheck: %d\n", config.HealthCheck)
	utils.Print(data.Green, "[+] Logs: %s\n", config.Logs)

	utils.Print(data.Green, "[+] Prohibited: \n")
	for _, dir := range config.Prohibited {
		utils.Print(data.Gray, "\t- %s\n", dir)
	}
}

// main function, call functions to read json, setup configurations, print errors and start the server
func main() {
	var config data.Config
	var servers server.ServersData

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

	utils.Print(data.Gray, "\nPress CTRL^C to stop\n")
	initializeServers(config.Servers, servers)
	requests.StartListener(config, servers)
}
