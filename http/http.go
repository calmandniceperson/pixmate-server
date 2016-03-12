package http

import (
	"imgturtle/io"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"os"
)

func Start() {
	r := mux.NewRouter().StrictSlash(true)
	r.NotFoundHandler = http.HandlerFunc(errorHandler)

	/*
	 *
	 *  ROUTES
	 */
	r.HandleFunc("/", mainPageHandler)
	r.HandleFunc("/upload", uploadHandler)
	r.HandleFunc("/apps", appsPageHandler)
	r.HandleFunc("/favicon.ico", favIcoHandler)
	r.HandleFunc("/img/{id}", imageHandler)
	r.HandleFunc("/{id}", imagePageHandler)
	r.HandleFunc("/api/upload", apiUploadHandler)
	r.HandleFunc("/error", errorHandler)

	n := negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(MiddleWare),
		negroni.NewLogger(),
		negroni.NewStatic(http.Dir(os.Getenv("APP_LOCATION"))),
		negroni.NewStatic(http.Dir("public")),
	)
	n.UseHandler(r)

	go func(n http.Handler) {
		cio.PrintMessage(0, "https: Running on port 8001...")
		err := http.ListenAndServeTLS(":8001", "http/ssl/cert.pem", "http/ssl/key.pem", n)
		if err != nil {
			log.Fatal(err)
		}
	}(n)

	cio.PrintMessage(0, "http: Running on port 8000...")
	http.ListenAndServe(":8000", n)
}
