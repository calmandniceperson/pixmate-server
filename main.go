/*
 * IMGTURTLE
 * GO PROTOTYPE
 * 2015
 */

package main

import (
	"imgturtle/http"
	"imgturtle/db"
)

func main() {
	db.Init()
	http.Init()
}
