package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"data"
	"utils"
)

func printJsonData(config data.Config) {
	utils.Print(data.Green, "[+] Using algorithm ", config.Algorithm)
	
	utils.Print(data.Green, "[+] Servers: ")
	for _, server := range config.Servers {
		fmt.Printf(" - %s\n", server)
	}

	if config.HealtCheck {
		utils.Print(data.Green, "[+] HealtCheck ", config.HealtCheck)
	} else {
		utils.Print(data.Red, "[+] HealtCheck ", config.HealtCheck)
	}
}

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
	printJsonData(config)
	
	
}