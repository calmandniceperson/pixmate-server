package http

import (
	"encoding/json"
	"imgturtle/db"
	"imgturtle/misc"
	"net/http"
	"path"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(cookieKey))

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
			misc.PrintMessage(1, "http", "requsr.go", "mePageHandler()", "500. Couldn't parse template.\n"+err.Error())
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		// return the template or print an error if one occurs
		if err := tmpl.Execute(w, user); err != nil {
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

type SignInData struct {
	Ue  string
	Pwd string
}

func signInHandler(w http.ResponseWriter, req *http.Request) {
	// If the request is a GET request
	// return the signin/signup page
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
		valid, uid, err := db.CheckUserCredentials(s.Ue, s.Pwd)
		if valid {
			// Get a session. We're ignoring the error resulted from decoding an
			// existing session: Get() always returns a session, even if empty.
			session, err := store.Get(req, "imgturtle")
			// Set some session values.
			session.Values["uid"] = uid
			session.Values["uname"] = s.Ue
			// Save it before we write to the response/return from the handler.
			session.Save(req, w)
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
