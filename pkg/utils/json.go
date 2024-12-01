package utils

import (
	"encoding/json"
	"io"
	"os"

	"data"
)

// reads the json and conver it to a better structure
func ReadAndParseJson(path string) data.Config {
	var rawConfig data.ConfigRaw
	
	err := ReadJson(&rawConfig, path)
		if err != nil {
			Print(data.Red, err.Error())
		}

	return rawConfig.Cast()
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
