/*
 * IMGTURTLE
 * GO PROTOTYPE
 * 2015
 */

package main

import (
	"imgturtle/db"
	"imgturtle/fs"
	"imgturtle/http"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	db.Start()

	if runtime.GOOS == "windows" {
		cmd := exec.Command("cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	// Check if the image storage exists before using it
	// if it doesn't, this function will create it
	fs.Start()

	http.Start()
}
