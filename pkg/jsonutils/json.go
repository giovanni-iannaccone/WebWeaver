package jsonutils

import (
	"encoding/json"
	"io"
	"os"

	"console"
	"config"
)

// reads the json and conver it to a better structure
func ReadAndParseJson(path string) config.Config {
	var rawConfig config.ConfigRaw
	var config config.Config

	err := ReadJson(&rawConfig, path)
	if err != nil {
		console.Println(console.Red, "%s %s", err.Error(), path)
	}

	config = rawConfig.Cast()
	config.Path = path
	return config
}

// reads a json file
func ReadJson(data any, path string) error {
	jsonFile, err := os.Open(path)

	if err != nil {
		return err
	}

	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &data)
	return nil
}
