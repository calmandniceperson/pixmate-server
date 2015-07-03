package http

import (
  "html/template"
  "os"
  "net/http"
  "path"
  "github.com/fatih/color"
  "github.com/gorilla/mux"
)

func MiddleWare(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc){
  if true{
    next(rw, r)
  }
}

func mainPageHandler(w http.ResponseWriter, req *http.Request) {
  color.Cyan("REQ: request for %s.html", req.URL.Path)
	resourcePath := "public" + req.URL.Path + ".html"
	if req.URL.Path == "/" {
	   resourcePath = "public/" + "welcome.html"
	}

  if _, err := os.Stat(resourcePath); os.IsNotExist(err) {
     color.Red("ERR: 404. no such file or directory: %s.html", req.URL.Path)
     http.Error(w, ("error 404. " + req.URL.Path + ".html could not be found."), 404)
  }else{
	   color.Green("INF: serving static file => %s", resourcePath)
	   http.ServeFile(w, req, resourcePath)
  }
}

type User struct{
  Uname string
}

func mePageHandler(w http.ResponseWriter, req *http.Request) {
  user := User{"Your "}

  fp := path.Join("public", "people.html")

  // form template
  tmpl, err := template.ParseFiles(fp)

  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  // return the template or print an error if one occurs
  if err := tmpl.Execute(w, user); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func peoplePageHandler(w http.ResponseWriter, req *http.Request){
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
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  // return the template or print an error if one occurs
  if err := tmpl.Execute(w, user); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

/*func imageHandler(w http.ResponseWriter, req *http.Request){
  vars := mux.Vars(req)
  id := vars["id"]
  resourcePath := "/public/img/" + id + ".jpg"

  if _, err := os.Stat(resourcePath); os.IsNotExist(err) {
    resourcePath := "/public/img/" + id + ".png"
    if _, err := os.Stat(resourcePath); os.IsNotExist(err) {
      color.Red("error 404. Image %s could not be found.", resourcePath)
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }else{
  	   color.Green("INF: serving static img => %s", resourcePath)
  	   http.ServeFile(w, req, resourcePath)
    }
  }else{
	   color.Green("INF: serving static img => %s", resourcePath)
	   http.ServeFile(w, req, resourcePath)
  }
}*/
