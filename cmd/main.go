package main

import (
	"flag"
	"os"

	"config"
	"console"
	"healthCheck"
	"jsonutils"
	"requests"
	"webui"
)

// calls the function to check errors, if any, print them
func checkAndPrintErrors(config config.Config) bool {
	errors := config.CheckValidity()
	if errors != nil {
		for _, err := range errors {
			console.Print(console.Red, "%s\n", err)
		}

		return true
	}

	return false
}

// merges the configurations from the file and those from cli
func mergeConfigs(cli *config.Config, file config.Config) {
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
func parseFlags(path *string, config *config.Config) {
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
	var config config.Config = *config.GetConfig()
	
	console.Println(console.Green, "[+] Using algorithm %s", config.Algorithm)
	console.Println(console.Green, "[+] Host: %s", config.Host)
	console.Println(console.Green, "[+] Dashboard: %d", config.Dashboard)

	console.Print(console.Green, "[+] Servers: \n")
	for i := range config.Servers.Inactive {
		console.Print(console.Gray, "\t- %s\n", config.Servers.Inactive[i])
	}

	console.Print(console.Green, "[+] HealthCheck: %d\n", config.HealthCheck)
	console.Print(console.Green, "[+] Logs: %s\n", config.Logs)

	console.Print(console.Green, "[+] Prohibited: \n")
	for _, dir := range config.Prohibited {
		console.Print(console.Gray, "\t- %s\n", dir)
	}
}

// prints help messages
func printHelp(_ string) error {
	var binaryName string = os.Args[0]

	console.Print(console.Reset, "%s\t\t--help\t\tShow this screen\n", binaryName)
	console.Print(console.Reset, "%s\t\t--config\tSpecify a configuration file\n", binaryName)
	console.Print(console.Reset, "( if the configuration isn't specified, the file will be ./configs/config.json )\n\n")
	console.Print(console.Reset, "Example: %s -c config.json\n\n", binaryName)
	console.Print(console.Reset, "Use a different value from the one in configurations by passing it as an arg\n")
	console.Print(console.Reset, "Example: %s --logs logs.txt", binaryName)
	
	os.Exit(0)
	return nil
}

// starts web ui and health check
func startNetworkStuff(config config.Config) {
	healthcheck.HealthCheck(config.Servers)
	
	if config.Dashboard >= 0 {
		go webui.RenderUI()
	}

	if config.HealthCheck >= 0 {
		go healthcheck.StartHealthCheckTimer(
			config.Servers, 
			config.HealthCheck, 
			config.Dashboard < 0,
		)
	}
}

// init gets executed before the main
func init() {
	webui.Init()
}

// main function, calls functions to read json, setup configurations, prints errors and starts the server
func main() {
	var configFilePath string

	var fileConfig config.Config
	var config *config.Config = config.GetConfig()

	parseFlags(&configFilePath, config)

	console.Println(console.Green, "===== Starting WebWeaver =====")
	console.Println(console.Reset, "Reading config files...")

	fileConfig = jsonutils.ReadAndParseJson(configFilePath)

	mergeConfigs(config, fileConfig)
	config.Path = configFilePath

	if checkAndPrintErrors(*config) {
		return
	}

	printConfigData()
	
	startNetworkStuff(*config)
	if (config.Dashboard >= 0) {
		console.Println(console.Blue, "Online, go to localhost:%d to access dashboard", config.Dashboard)
	}

	console.Print(console.Gray, "\nPress CTRL^C to stop\n")
	requests.StartListener()
}