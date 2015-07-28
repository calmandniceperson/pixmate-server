package http

import (
	"encoding/json"
	"imgturtle/db"
	"net/http"
	"path"
	"text/template"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

// User struct stores user data
// to fill into the user's profile page
type User struct {
	Uname string
}

func mePageHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		user := User{"Your "}

		fp := path.Join("public", "people.html")

		// form template
		tmpl, err := template.ParseFiles(fp)

		if err != nil {
			color.Red("ERR: 500. Couldn't parse template.\n%s", err.Error())
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		// return the template or print an error if one occurs
		if err := tmpl.Execute(w, user); err != nil {
			color.Red("ERR: 500. Couldn't return template.\n%s", err.Error())
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
		} else {
			color.Green("INF: serving static file => %s", "people.html (me)")
		}
	}
}

func peoplePageHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		// get variables (name in people/{name}) from request
		vars := mux.Vars(req)

		// get name
		name := vars["name"]

		// insert name into User object
		user := User{name + "'s "}

		fp := path.Join("public", "people.html")

		// form template
		tmpl, err := template.ParseFiles(fp)

		if err != nil {
			color.Red("ERR: 500. Couldn't parse template.")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// return the template or print an error if one occurs
		if err := tmpl.Execute(w, user); err != nil {
			color.Red("ERR: 500. Couldn't return template.")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			color.Green("INF: serving static file => %s", "people.html")
		}
	}
}

func signInHandler(w http.ResponseWriter, req *http.Request) {
	// If the request is a GET request
	// return the signin/signup page
	if req.Method == "GET" {
		color.Cyan("INF: serving static file => %s", "signin.html")
		http.ServeFile(w, req, "public/signin.html")
	} else if req.Method == "POST" { // If the request is a POST request sign the user in

	}
}

// SignUpData stores the
// data that is sent with a
// signup POST request
type SignUpData struct {
	Uname string
	Pwd   string
	Email string
}

func signUpHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		color.Cyan("INF: serving static file => %s", "signup.html")
		http.ServeFile(w, req, "public/signup.html")
	} else if req.Method == "POST" {
		decoder := json.NewDecoder(req.Body)
		var s SignUpData
		err := decoder.Decode(&s)
		if err != nil {
			color.Red("ERR: Error decoding signup JSON.\n%s", err.Error())
		}

		err = db.InsertNewUser(s.Uname, s.Pwd, s.Email)
		if err != nil {
			color.Red(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			//http.Redirect(w, req, "/", 200)
			http.StatusText(200)
		}
	}
}
