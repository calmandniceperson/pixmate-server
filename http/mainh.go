package http

import (
	"imgturtle/io"
	"net/http"
	"path"
)

func MiddleWare(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if true { // for development, later insert a real condition
		next(rw, r)
	}
}

func errorHandler(w http.ResponseWriter, req *http.Request) {
	cio.PrintMessage(1, req.URL.Path+" not found. Serving error.html")
	http.ServeFile(w, req, "public/error.html")
}

func mainPageHandler(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "public/imgturtle.html")
}

func favIcoHandler(w http.ResponseWriter, req *http.Request) {
	fp := path.Join("public/img/", "favicon.ico")
	http.ServeFile(w, req, fp)
}
