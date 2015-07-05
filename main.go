/*
 * IMGTURTLE
 * GO PROTOTYPE
 * 2015
 */

package main

import (
	"imgturtle/http"
	"imgturtle/db"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	db.Init()

	if runtime.GOOS == "windows"{
		cmd := exec.Command("cls")
		cmd.Stdout = os.Stdout
	  cmd.Run()
	}else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
	  cmd.Run()
	}


	http.Init()
}
