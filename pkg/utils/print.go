package utils

import (
	"data"
	"fmt"
)

// acts exactly like println but print colored text
func Print(color string, text ...any) {
	fmt.Print(color)
	fmt.Print(text...)
	fmt.Println(data.Reset)
}
