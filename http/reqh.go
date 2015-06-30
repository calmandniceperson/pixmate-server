package http

import (
  "net/http"
  "github.com/fatih/color"
)

func mainPageHandler(w http.ResponseWriter, req *http.Request) {
	resourcePath := "public/html" + req.URL.Path + ".html"
	if req.URL.Path == "/" {
		resourcePath = "public/html/" + "welcome.html"
	}
	color.Red("serving static file => %s", resourcePath)
	http.ServeFile(w, req, resourcePath)
}
