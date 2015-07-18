package http

import (
	"crypto/rand"
	"errors"
	"fmt"
	"html/template"
	"imgturtle/fs"
	"io"
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

func uploadHandler(w http.ResponseWriter, req *http.Request) {
	// Check if the image storage exists before using it
	// if it doesn't, this function will create it
	// This might later be moved to main.go if needed
	fs.CreateImageStorageIfNotExists()

	if req.Method == "POST" {
		req.ParseMultipartForm(32 << 20)
		file, fileheader, err := req.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", fileheader.Header)

		filename, err := findAvailableName(strings.Split(fileheader.Filename, ".")[1])
		if err != nil {
			color.Red(err.Error())
			return
		}

		if filename != "" {
			filename = fs.ImgStoragePath + filename[0:2] + "/" + filename[2:4] + "/" + filename[4:6] + "/" + filename[6:8] + "/" + filename[7:len(filename)] + "." + strings.Split(fileheader.Filename, ".")[1]
			f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
			defer f.Close()
			if err != nil {
				color.Red(err.Error())
				return
			}
			color.Green("INF: File " + filename + " has been created.")
			bytesCopied, err := io.Copy(f, file)
			if err != nil {
				color.Red(err.Error())
				return
			}
			color.Green("INF: Content of uploaded image (" + strconv.FormatInt(bytesCopied, 10) + " Bytes) has been copied to " + filename + ".")
		}
	}
}

func findAvailableName(fileExt string) (string, error) {
	b := make([]byte, 16)
	rand.Read(b)
	filename := fmt.Sprintf("%x", b)

	// check if the file already exists
	// to instantly generate a new name
	// before running through the whole process
	if _, err := os.Stat(fs.ImgStoragePath + filename[0:2] + "/" + filename[2:4] + "/" + filename[4:6] + "/" + filename[6:8] + "/" + filename[9:len(filename)] + "." + fileExt); os.IsNotExist(err) {
		// doesn't exist

		// OFFSET
		// A value to set where
		// to start slicing the string
		offset := 0

		// CURRENTPATH
		// Stores the path the
		// loop has already worked through
		currentPath := ""

		// Endless, since the return or the
		// recursive function call will end
		// the function later on anyway
		for true {

			// If the OFFSET is lower than 8
			// there are still directories to be created
			// since there need to be at least 4 directories
			if offset < 8 {

				// Check if the directory that is to be created
				// already exists.
				//
				// If it doesn't the directory will be created,
				// the directory will be added to the current path
				// and the offset will be increased by the number of
				// letters used (2) from the filename random hash
				if _, err := os.Stat(fs.ImgStoragePath + currentPath + filename[offset:offset+2] + "/"); os.IsNotExist(err) {
					// doesn't exist
					color.Cyan("INF: ../" + currentPath + filename[offset:offset+2] + "/ created.")
					os.Mkdir(fs.ImgStoragePath+currentPath+filename[offset:offset+2]+"/", 0776)
					currentPath += filename[offset:offset+2] + "/"
					offset += 2
				} else {
					// exists
					color.Cyan("INF: ../" + currentPath + filename[offset:offset+2] + "/ already existed and thus is not created!")
					currentPath += filename[offset:offset+2] + "/"
					offset += 2
				}

				// If the OFFSET is higher than 9
				// the next step is to create the file itself.
				// The function won't do this itself.
				// Instead it will check, if the file exists (again, just to be sure)
				// and if the file doesn't exist, it will return the name of the file
				// and an empty error.
			} else {
				if _, err := os.Stat(fs.ImgStoragePath + currentPath + filename[offset:len(filename)] + "." + fileExt); os.IsNotExist(err) {
					// doesn't exist
					return filename, nil
				}
			}
		}
	} else {
		// The file already exists
		// so a new name has to be generated
		findAvailableName(fileExt)
	}
	return "", errors.New("Something went wrong. (findAvailableName)")
}
