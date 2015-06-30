package http

import (
  "net/http"
  "github.com/fatih/color"
)

func Init(){
  color.Green("IMGCAT Server running on port 8000")
	http.HandleFunc("/", mainPageHandler)
	http.ListenAndServe(":8000", nil)
}
