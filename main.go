/*
 * IMGTURTLE
 * GO PROTOTYPE
 * 2015
 */

package main

import (
	"os"
	"os/exec"
	"pixmate-server/db"
	"pixmate-server/fsys"
	"pixmate-server/http"
	"runtime"
)

func main() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	db.Start()
	// Check if the image storage exists before using it
	// if it doesn't, this function will create it
	fsys.Start()
	go func() {
		fsys.RemoveOldImages()
	}()
	http.Start()
}
