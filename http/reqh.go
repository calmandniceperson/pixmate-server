package http

import (
	"html/template"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

// MiddleWare describes a process (like checking for a valid user id)
// on every request it is being used on
func MiddleWare(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if true { // for the time being in development
		next(rw, r)
	}
}

func mainPageHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		color.Cyan("REQ: request for %s.html", req.URL.Path)
		resourcePath := "public" + req.URL.Path + ".html"
		if req.URL.Path == "/" {
			resourcePath = "public/" + "welcome.html"
		}

		if _, err := os.Stat(resourcePath); os.IsNotExist(err) {
			color.Red("ERR: 404. no such file or directory: %s.html", req.URL.Path)
			http.Error(w, ("error 404. " + req.URL.Path + ".html could not be found."), 404)
		} else {
			color.Green("INF: serving static file => %s", resourcePath)
			http.ServeFile(w, req, resourcePath)
		}
	}
}

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
			color.Red("ERR: 500. Couldn't parse template.")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// return the template or print an error if one occurs
		if err := tmpl.Execute(w, user); err != nil {
			color.Red("ERR: 500. Couldn't return template.")
			http.Error(w, err.Error(), http.StatusInternalServerError)
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

// Img struct stores image data
// in order to display the correct image and its data
// in the image page
type Img struct {
	ImgTitle    string
	ImgFilePath string
}

func imageHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		/*
		 * fetch image id from URL parameters (e.g. /img/123)
		 */
		vars := mux.Vars(req)
		id := vars["id"]

		/*
		 * 1. set resource path to directory containing images
		 * 2. check if the directory exists (if not, print an error)
		 */
		resourcePath := "public/img/"
		d, err := os.Open(resourcePath)
		if err != nil {
			color.Red(err.Error())
		}

		/*
		 * close the file system connection
		 *
		 * from GO documentation:
		 * "Defer is used to ensure that a function call is performed
		 * later in a program’s execution, usually for purposes of cleanup.
		 * defer is often used where e.g. ensure and finally would be used
		 * in other languages."
		 */
		defer d.Close() // this will be executed at the end of the enclosing function

		/*
		 * Read file info
		 *
		 * "Readdir reads the contents of the directory associated
		 * with file and returns a slice of up to n FileInfo values"
		 */
		fi, err := d.Readdir(-1)

		if err != nil {
			color.Red(err.Error())
		}

		/*
		 * Iterate through the files in /public/img
		 * and try to find a fitting image (same name, file extension, etc.)
		 */

		matches := 0 // match count
		for _, fi := range fi {
			if matches > 0 {
				return
			}
			if fi.Mode().IsRegular() { // if there are no mode type bits set
				if strings.Split(fi.Name(), ".")[0] == id { // if the file name matches the given image ID
					img := Img{strings.Split(fi.Name(), ".")[0], (fi.Name())}
					fp /*file path*/ := path.Join("public", "img.html")
					tmpl, err := template.ParseFiles(fp)

					if err != nil {
						color.Red("ERR: 500. Couldn't parse template.")
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					// return the template or print an error if one occurs
					if err := tmpl.Execute(w, img); err != nil {
						color.Red("ERR: 500. Couldn't return template.")
						http.Error(w, err.Error(), http.StatusInternalServerError)
					} else {
						color.Green("INF: serving static file => %s with image %s (size: %s Bytes)", "img.html", fi.Name(), strconv.FormatInt(fi.Size(), 10))
						matches++
					}
				}
			}
		}

		/*
		 * if no images were found, return text (for now, maybe HTML later)
		 */
		if matches == 0 {
			w.Write([]byte("Sorry. We couldn't find an image called " + id + "."))
		}
	}
}
