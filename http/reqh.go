package http

import (
	"html/template"
	"net/http"
	"path"

	"github.com/fatih/color"
)

// MiddleWare describes a process (like checking for a valid user id)
// on every request it is being used on
func MiddleWare(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if true { // for the time being in development
		next(rw, r)
	}
}

// Page contains
// the data that is inserted
// into the welcome page template
type Page struct {
	IsLoggedIn bool
}

func errorHandler(w http.ResponseWriter, req *http.Request) {
	color.Red("INF: request for %s. 404, not found.\nINF: serving static file => %s", req.URL.Path, "error.html")
	http.ServeFile(w, req, "public/error.html")
}

func mainPageHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		pdata := Page{false}
		fp /*file path*/ := path.Join("public", "welcome.html")
		// parse img.html as template
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			color.Red("ERR: 500. Couldn't parse template\n%s", err.Error())
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		// return (execute) the template or print an error if one occurs
		if err := tmpl.Execute(w, pdata); err != nil {
			color.Red("ERR: 500. Couldn't return template.\n%s", err.Error())
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
		} else {
			color.Green("INF: serving static file => %s", "welcome.html")
		}
	}
}
