package misc

import (
	"github.com/fatih/color"
)

func PrintMessage(cl int, pkg string, file string, funct string, msg string) {
	if cl == 0 {
		color.Green("%s: %s %s -> %s", pkg, file, funct, msg)
	} else if cl == 1 {
		color.Red("%s: %s %s -> %s", pkg, file, funct, msg)
	} else if cl == 2 {
		color.Cyan("%s: %s %s -> %s", pkg, file, funct, msg)
	}
}
