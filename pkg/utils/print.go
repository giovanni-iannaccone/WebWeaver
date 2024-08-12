package utils

import (
	"data"
	"fmt"
)

// acts exactly like println but print colored text
func Print(color string, format string, args ...any) {
	fmt.Print(color)
	fmt.Printf(format, args...)
	fmt.Print(data.Reset)
}
