/*
 * IMGTURLTE
 * GO PROTOTYPE
 * 2015
 */

package http

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

// Start is the http packages launch method
// It creates all the routes, adds negroni logging and starts the server
func Start() {
	color.Cyan("http: Starting HTTP/HTTPS server...")
	r := mux.NewRouter().StrictSlash(true)
	r.NotFoundHandler = http.HandlerFunc(errorHandler)

	// main page
	r.HandleFunc("/", mainPageHandler)

	// custome profile URL
	r.HandleFunc("/u/{name}", peoplePageHandler)

	// personal profile url
	r.HandleFunc("/me", mePageHandler)

	// Sign in GET & POST
	r.HandleFunc("/signin", signInHandler)

	r.HandleFunc("/signup", signUpHandler)

	// file upload
	r.HandleFunc("/upload", uploadHandler)

	// errors
	r.HandleFunc("/error", errorHandler)

	r.HandleFunc("/favicon.ico", favIcoHandler)

	// images
	r.HandleFunc("/img/{id}", imageHandler)
	r.HandleFunc("/{id}", imagePageHandler)

	// initialise negroni
	// include middleware, logger, etc.
	n := negroni.New(
		negroni.NewRecovery(),
		//negroni.HandlerFunc(MiddleWare),
		//negroni.NewLogger(),
		negroni.NewStatic(http.Dir("public")),
	)
	n.UseHandler(r)
	color.Green("http: Running on port 8000")

	http.ListenAndServe(":8000", n)
	//http.ListenAndServeTLS(port, certificate.pem, key.pem, nil) for https
}
