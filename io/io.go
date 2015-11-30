package cio

import (
	"github.com/fatih/color"
)

func PrintMessage(cl int, msg string) {
	if cl == 0 {
		color.Green("%s", msg)
	} else if cl == 1 {
		color.Red("%s", msg)
	} else if cl == 2 {
		color.Cyan("%s", msg)
	}
}
