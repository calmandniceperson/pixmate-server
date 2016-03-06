package cio

import (
	"github.com/fatih/color"
)

func PrintMessage(cl int, msg string) {
	if cl == 0 {
		color.Green("[pixmate] %s", msg)
	} else if cl == 1 {
		color.Red("[pixmate] %s", msg)
	} else if cl == 2 {
		color.Cyan("[pixmate] %s", msg)
	}
}
