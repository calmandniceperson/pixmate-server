package http

import (
	"encoding/json"
	"imgturtle/db"
	"imgturtle/misc"
	"net/http"
	"path"
	"text/template"

	"github.com/gorilla/mux"
)

// User struct stores user data
// to fill into the user's profile page
type User struct {
	Uname string
}

// UPageData stores data for the people
// page
type UPageData struct {
	IsLoggedIn   bool
	IsFollowable bool
	User         User
}

func mePageHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		var udata UPageData
		if uName, err := getUserCookieData(req); err == nil {
			udata.IsLoggedIn = true
			udata.User.Uname = uName
		} else {
			udata.IsLoggedIn = false
		}
		udata.IsFollowable = false
		fp := path.Join("public", "people.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			misc.PrintMessage(1, "http", "requsr.go", "mePageHandler()", "500. Couldn't parse template.\n"+err.Error())
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		// return the template or print an error if one occurs
		if err := tmpl.Execute(w, udata); err != nil {
			misc.PrintMessage(1, "http", "requsr.go", "mePageHandler()", "500. Couldn't return template.\n"+err.Error())
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
		} else {
			misc.PrintMessage(2, "http", "requsr.go", "mePageHandler()", "serving file people.html(me)")
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
			misc.PrintMessage(1, "http", "requsr.go", "peoplePageHandler()", "500. Couldn't parse template.")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// return the template or print an error if one occurs
		if err := tmpl.Execute(w, user); err != nil {
			misc.PrintMessage(1, "http", "requsr.go", "peoplePageHandler()", "500. Couldn't return template.")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			misc.PrintMessage(2, "http", "requsr.go", "peoplePageHandler()", "serving file people.html")
		}
	}
}

// SignInData holds all the sign-in related
// data during the signin process
type SignInData struct {
	Ue  string
	Pwd string
}

func signInHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		misc.PrintMessage(2, "http", "requsr.go", "signInHandler()", "serving file signin.html")
		http.ServeFile(w, req, "public/signin.html")
	} else if req.Method == "POST" { // If the request is a POST request sign the user in
		decoder := json.NewDecoder(req.Body)
		var s SignInData
		err := decoder.Decode(&s)
		if err != nil {
			misc.PrintMessage(1, "http", "requsr.go", "signInHandler()", "Error decoding signup JSON\n"+err.Error())
		}
		valid, err := db.CheckUserCredentials(s.Ue, s.Pwd)
		if valid {
			setUserCookie(s.Ue, w)
			http.StatusText(200)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		} else {
			if err != nil {
				misc.PrintMessage(1, "http", "requsr.go", "signInHandler()", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
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
		misc.PrintMessage(2, "http", "requsr.go", "signUpHandler()", "serving file signup.html")
		http.ServeFile(w, req, "public/signup.html")
	} else if req.Method == "POST" {
		decoder := json.NewDecoder(req.Body)
		var s SignUpData
		err := decoder.Decode(&s)
		if err != nil {
			misc.PrintMessage(1, "http", "requsr.go", "signUpHandler()", "Error decoding signup JSON.\n"+err.Error())
		}

		err = db.InsertNewUser(s.Uname, s.Pwd, s.Email)
		if err != nil {
			misc.PrintMessage(1, "http", "requsr.go", "signUpHandler()", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.StatusText(200)
		}
	}
}

func logoutHandler(w http.ResponseWriter, req *http.Request) {
	clearUserCookie(w)
	http.Redirect(w, req, "/", 200)
}
