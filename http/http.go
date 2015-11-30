package http

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

func Start() {
	r := mux.NewRouter().StrictSlash(true)
	r.NotFoundHandler = http.HandlerFunc(errorHandler)

	/*
	 *
	 *  ROUTES
	 *
	 */
	r.HandleFunc("/", mainPageHandler)
	r.HandleFunc("/upload", uploadHandler)
	r.HandleFunc("/error", errorHandler)
	r.HandleFunc("/favicon.ico", favIcoHandler)
	r.HandleFunc("/img/{id}", imageHandler)
	r.HandleFunc("/{id}", imagePageHandler)

	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewStatic(http.Dir("public")),
	)
	n.UseHandler(r)

	color.Green("http: Running on port 8000...")
	http.ListenAndServe(":8000", n)
}
