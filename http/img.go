package http

import (
	"crypto/rand"
	"errors"
	"fmt"
	"imgturtle/db"
	"imgturtle/fsys"
	"imgturtle/io"
	"net/http"
	"path"
	"strconv"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

type Img struct {
	ImgTitle    string
	ImgFilePath string
}

func imagePageHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		vars := mux.Vars(req)
		id := vars["id"]
		if len(id) > fsys.ImgNameLength {
			if strings.Contains(id, ".") {
				id = strings.Split(id, ".")[0]
			}
			found, imgPath, imgID, title, errc, err := db.CheckIfImageExists(id)
			if err != nil {
				cio.PrintMessage(1, err.Error())
				if errc == http.StatusInternalServerError {
					http.Error(w, err.Error(), errc)
					return
				} else if errc == http.StatusNotFound {
					errorHandler(w, req)
					return
				}
			}
			if found == true {
				img := Img{title, "/img/" + imgID}
				fp := path.Join("public", "img.html")
				tmpl, err := template.ParseFiles(fp)
				if err != nil {
					cio.PrintMessage(1, err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if err := tmpl.Execute(w, img); err != nil {
					cio.PrintMessage(1, err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				cio.PrintMessage(0, imgPath+" rendered into img.html")
				return
			}
			errorHandler(w, req)
			return
		}
		errorHandler(w, req)
		return
	}
}

func imageHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		vars := mux.Vars(req)
		id := vars["id"]
		if len(id) > fsys.ImgNameLength {
			if strings.Contains(id, ".") {
				id = strings.Split(id, ".")[0]
			}
			found, imgPath, _, _, errc, err := db.CheckIfImageExists(id)
			if err != nil {
				cio.PrintMessage(1, (string(errc) + " " + err.Error()))
				if errc == http.StatusInternalServerError {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				} else if errc == http.StatusNotFound {
					errorHandler(w, req)
					return
				}
			}
			if found == true {
				fp := path.Join(fsys.ImgStoragePath, imgPath)
				http.ServeFile(w, req, fp)
				return
			}
			http.Error(w, errors.New("Image with ID "+id+" could not be found.").Error(), http.StatusNotFound)
			return
		}
		http.Error(w, errors.New(id+" is not a valid ID.").Error(), http.StatusNotFound)
		return
	}
}

func uploadHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		req.ParseMultipartForm(32 << 20)
		file, fileheader, err := req.FormFile("uploadfile")
		if err != nil {
			cio.PrintMessage(1, err.Error())
			return
		}
		defer file.Close()
		var id string
		created := false
		for created != true {
			id = generateImageID()
			err = db.CheckIfImageIDInUse(id)
			if err != nil {
				if err.Error() == "Image ID '"+id+"' exists.'" {
					created = false
					continue
				} else {
					created = false
					cio.PrintMessage(1, err.Error())
					return
				}
			}
			created = true
		}
		filePath := fsys.ImgStoragePath + id + "." + strings.Split(fileheader.Filename, ".")[1]
		err = db.StoreImage(id, strings.Split(fileheader.Filename, ".")[0] /*image name*/, filePath, strings.Split(fileheader.Filename, ".")[1] /*extension*/)
		if err != nil {
			cio.PrintMessage(1, err.Error())
			return
		}
		bytesCopied, err := fsys.StoreImage(filePath, file)
		if err != nil {
			cio.PrintMessage(1, err.Error())
			return
		}
		cio.PrintMessage(0, "File "+filePath+" has been created.")
		if err != nil {
			cio.PrintMessage(1, err.Error())
			return
		}
		cio.PrintMessage(0, "Content of uploaded image ("+strconv.FormatInt(bytesCopied, 10)+" Bytes) has been copied to "+filePath+".")
		http.Redirect(w, req, "/"+id, http.StatusFound)
	}
}

func generateImageID() string {
	// generate a random hash
	// length is provided by an environment
	// variable or user input
	b := make([]byte, fsys.ImgNameLength)
	rand.Read(b)
	filename := fmt.Sprintf("%x", b)
	return filename
}

