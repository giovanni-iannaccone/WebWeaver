package main

import (
	"os"
	"strconv"

	"data"
	"internals/requests"
	"utils"
	"webui"
)

var config data.Config

// calls the function to check errors, if any, print them
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

// merges the configurations from the file and those from cli
func mergeConfigs(cli *data.Config, file data.Config) {
	if cli.Algorithm == "" && file.Algorithm != cli.Algorithm {
		cli.Algorithm = file.Algorithm
	}

	if cli.Host == "" && file.Host != cli.Host {
		cli.Host = file.Host
	} 

	if cli.Dashboard < 0 && file.Dashboard != cli.Dashboard {
		cli.Dashboard = file.Dashboard
	} 

	if cli.HealthCheck < 0 && file.HealthCheck != cli.HealthCheck {
		cli.HealthCheck = file.HealthCheck
	}

	if cli.Logs == "" && file.Logs != cli.Logs {
		cli.Logs = file.Logs
	} 

	if cli.Prohibited == nil && len(file.Prohibited) > 0 {
		cli.Prohibited = file.Prohibited
	} 

	if cli.Servers == nil && len(file.Servers) > 0 {
		cli.Servers = file.Servers
	} 
}

// prints help and returns the configurations file path based on command line args
func parseArguments(args []string, config *data.Config) string {
	for i := range args {
		if args[i] == "--help" || args[i] == "-h" {
			printHelp(args)
			os.Exit(0)

		} else if args[i] == "--config" || args[i] == "-c" {
			return args[i + 1]

		} else if args[i] == "--algorithm" {
			config.Algorithm = args[i + 1]

		} else if args[i] == "--host" {
			config.Host = args[i + 1]
		
		} else if args[i] == "--dashboard" {
			config.Dashboard, _ = strconv.Atoi(args[i + 1])

		} else if args[i] == "--healthcheck" {
			config.HealthCheck, _ = strconv.Atoi(args[i + 1])

		} else if args[i] == "--logs" {
			config.Logs = args[i + 1]
		}
	}

	return "./configs/config.json"
}

// prints configurations
func printConfigData(config data.Config) {
	utils.Print(data.Green, "[+] Using algorithm %s\n", config.Algorithm)
	utils.Print(data.Green, "[+] Host: %s\n", config.Host)
	utils.Print(data.Green, "[+] Dashboard: %d\n", config.Dashboard)

	utils.Print(data.Green, "[+] Servers: \n")
	for _, server := range config.Servers {
		utils.Print(data.Gray, "\t- %s\n", server.URL.String())
	}

	utils.Print(data.Green, "[+] HealthCheck: %d\n", config.HealthCheck)
	utils.Print(data.Green, "[+] Logs: %s\n", config.Logs)

	utils.Print(data.Green, "[+] Prohibited: \n")
	for _, dir := range config.Prohibited {
		utils.Print(data.Gray, "\t- %s\n", dir)
	}
}

// prints help messages
func printHelp(args []string) {
	utils.Print(data.Reset, "%s\t\t--help\t | -h\t\tShow this screen\n", args[0])
	utils.Print(data.Reset, "%s\t\t--config | -c\t\t Specify a configuration file\n", args[0])
	utils.Print(data.Reset, "( if the configuration isn't specified, the file will be configs/config.json )\n\n")
	utils.Print(data.Reset, "Example: %s -c config.json\n\n", args[0])
	utils.Print(data.Reset, "Use a different value from the one in configurations by passing it as an arg\n")
	utils.Print(data.Reset, "Example: %s --logs logs.txt", args[0])
}

// init gets executed before the main
func init() {
	config = data.Config{
		Algorithm: "", 
		Host: "", 
		Dashboard: -1,
		Servers: nil,
		HealthCheck: -1,
		Logs: "",
		Prohibited: nil,
	}

	webui.Init()
}

// main function, calls functions to read json, setup configurations, prints errors and starts the server
func main() {
	var fileConfig data.Config

	var configFilePath string = parseArguments(os.Args, &config)

	utils.Print(data.Green, "===== Starting WebWeaver =====\n")
	utils.Print(data.Green, "[+] Reading config files\n")

	fileConfig = utils.ReadAndParseJson(configFilePath)

	mergeConfigs(&config, fileConfig)
	config.Path = configFilePath

	if checkAndPrintErrors(config) {
		return
	}

	printConfigData(config)

	webui.RenderUI(&config)
	utils.Print(data.Blue, "Online, go to localhost:%d to access dashboard", config.Dashboard)

	utils.Print(data.Gray, "\nPress CTRL^C to stop\n")
	requests.StartListener(&config)
}