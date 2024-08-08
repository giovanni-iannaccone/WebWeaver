package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"data"
	"utils"
)

// Checks for configuration validity: a valid algorithm and Server field not empty
func checkConfigValidity(config data.Config) []string {
	var errors []string

	if _, exists := data.Algorithms[config.Algorithm]; !exists {
		errors = append(errors, "[+] Unable to find algorithm " + config.Algorithm)
	}

	if len(config.Servers) == 0 {
		errors = append(errors, "[+] No server found ")
	}

	return errors
}

// if the configurations are valid, print them
func printJsonData(config data.Config) {
	utils.Print(data.Green, "[+] Using algorithm ", config.Algorithm)
	
	utils.Print(data.Green, "[+] Servers: ")
	for _, server := range config.Servers {
		fmt.Printf("\t- %s\n", server)
	}
		
	utils.Print(data.Green, "[+] HealtCheck: ", config.HealtCheck)
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

	utils.Print(data.Green, "===== Starting SilkRoute =====\n")
	
	utils.Print(data.Green, "[+] Reading config files")
	readJson(&config)

	errors := checkConfigValidity(config)
	if  errors != nil {
		for _, err := range errors {
			utils.Print(data.Red, err)
		}

		return 
	}

	printJsonData(config)

	
}