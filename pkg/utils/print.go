package utils

import (
	"fmt"
	"data"
)

func Print(color string, text ...any) {
	fmt.Print(color)
	fmt.Print(text...)
	fmt.Println(data.Reset)
}