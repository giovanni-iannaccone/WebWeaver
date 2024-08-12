package utils

import (
	"log"
	"os"
)

// Write data to file ( creates it if it doesn't exist, otherwise append)
func WriteLogs(data string, filePath string) {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Print(err)
	}

	f.WriteString(data)
}
