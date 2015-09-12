package http

import (
	"html/template"
	"imgturtle/misc"
	"net/http"
	"path"
)

// MiddleWare describes a process (like checking for a valid user id)
// on every request it is being used on
func MiddleWare(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if true { // for the time being in development
		next(rw, r)
	}
}

// MainPageData contains
// the data that is inserted
// into the main page template
type MainPageData struct {
	IsLoggedIn bool
}

func errorHandler(w http.ResponseWriter, req *http.Request) {
	misc.PrintMessage(1, "http", "reqh.go", "errorHandler()", req.URL.Path+" not found. Serving file error.html")
	http.ServeFile(w, req, "public/error.html")
}

func mainPageHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		var mdata MainPageData
		if _, err := getUserCookieData(req); err == nil {
			mdata.IsLoggedIn = true
		} else {
			mdata.IsLoggedIn = false
		}
		fp := path.Join("public", "imgturtle.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			misc.PrintMessage(1, "http", "reqh.go", "mainPageHandler()", "Couldn't parse template\n"+err.Error())
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		// return (execute) the template or print an error if one occurs
		if err := tmpl.Execute(w, mdata); err != nil {
			misc.PrintMessage(1, "http", "reqh.go", "mainPageHandler()", "Couldn't return template\n"+err.Error())
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
		} else {
			misc.PrintMessage(0, "http", "reqh.go", "mainPageHandler()", "serving file welcome.html")
		}
	}
}
