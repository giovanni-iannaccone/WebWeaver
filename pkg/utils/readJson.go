package utils

import (
	"encoding/json"
	"io"
	"os"
)

// read a json file
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
