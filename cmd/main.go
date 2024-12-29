package main

import (
	"flag"
	"os"

	"data"
	"internals/requests"
	"utils"
	"webui"
)

// calls the function to check errors, if any, print them
func checkAndPrintErrors(config data.Config) bool {
	errors := config.CheckValidity()
	if errors != nil {
		for _, err := range errors {
			utils.Print(data.Red, "%s\n", err)
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

	cli.Prohibited = file.Prohibited
	cli.Servers = file.Servers
}

// parses flags using the flag package
func parseFlags(path *string, config *data.Config) {
	flag.StringVar(path, "config", "./configs/config.json", "Config file path")
	flag.StringVar(&config.Algorithm, "algorithm", "", "Algorithm we will use to send packages")
	flag.StringVar(&config.Host, "host", "", "Address load balancer will listen on")
	flag.StringVar(&config.Logs, "logs", "", "File where logs will be saved")

	flag.IntVar(&config.Dashboard, "dashboard", -1, "Port where to start the dashboard")
	flag.IntVar(&config.HealthCheck, "healthcheck", -1, "Healthcheck timer")

	flag.BoolFunc("help", "Show help screen", printHelp)
	flag.Parse()
}

// prints configurations
func printConfigData() {
	var config data.Config = *data.GetConfig()
	
	utils.Print(data.Green, "[+] Using algorithm %s\n", config.Algorithm)
	utils.Print(data.Green, "[+] Host: %s\n", config.Host)
	utils.Print(data.Green, "[+] Dashboard: %d\n", config.Dashboard)

	utils.Print(data.Green, "[+] Servers: \n")
	for i := range config.Servers.Inactive {
		utils.Print(data.Gray, "\t- %s\n", config.Servers.Inactive[i])
	}

	utils.Print(data.Green, "[+] HealthCheck: %d\n", config.HealthCheck)
	utils.Print(data.Green, "[+] Logs: %s\n", config.Logs)

	utils.Print(data.Green, "[+] Prohibited: \n")
	for _, dir := range config.Prohibited {
		utils.Print(data.Gray, "\t- %s\n", dir)
	}
}

// prints help messages
func printHelp(_ string) error {
	var binaryName string = os.Args[0]

	utils.Print(data.Reset, "%s\t\t--help\t\tShow this screen\n", binaryName)
	utils.Print(data.Reset, "%s\t\t--config\tSpecify a configuration file\n", binaryName)
	utils.Print(data.Reset, "( if the configuration isn't specified, the file will be ./configs/config.json )\n\n")
	utils.Print(data.Reset, "Example: %s -c config.json\n\n", binaryName)
	utils.Print(data.Reset, "Use a different value from the one in configurations by passing it as an arg\n")
	utils.Print(data.Reset, "Example: %s --logs logs.txt", binaryName)
	
	os.Exit(0)
	return nil
}

// init gets executed before the main
func init() {
	webui.Init()
}

// main function, calls functions to read json, setup configurations, prints errors and starts the server
func main() {
	var configFilePath string

	var fileConfig data.Config
	var config *data.Config = data.GetConfig()

	parseFlags(&configFilePath, config)

	utils.Print(data.Green, "===== Starting WebWeaver =====\n")
	utils.Print(data.Green, "[+] Reading config files\n")

	fileConfig = utils.ReadAndParseJson(configFilePath)

	mergeConfigs(config, fileConfig)
	config.Path = configFilePath

	if checkAndPrintErrors(*config) {
		return
	}

	printConfigData()

	if (config.Dashboard >= 0) {
		go webui.RenderUI()
		utils.Print(data.Blue, "Online, go to localhost:%d to access dashboard", config.Dashboard)
	}
	
	utils.Print(data.Gray, "\nPress CTRL^C to stop\n")
	requests.StartListener()
}