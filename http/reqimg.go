package http

import (
	"crypto/rand"
	"errors"
	"fmt"
	"imgturtle/db"
	"imgturtle/fs"
	"imgturtle/misc"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

// Img struct stores image data
// in order to display the correct image and its data
// in the image page
type Img struct {
	ImgTitle    string
	ImgFilePath string
}

func imageHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		// fetch image ID from url
		vars := mux.Vars(req)
		id := vars["id"]
		if len(id) > fs.ImgNameLength {
			if strings.Contains(id, ".") {
				id = strings.Split(id, ".")[0]
			}
			found, imgPath, _, _, errc, err := db.CheckIfImageExists(id)
			if err != nil {
				misc.PrintMessage(1, "http", "reqimg.go", "imageHandler()", (string(errc) + " " + err.Error()))
				if errc == 500 {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				} else if errc == 404 {
					errorHandler(w, req)
					return
				}
			}
			if found == true {
				fp := path.Join(fs.ImgStoragePath, imgPath)
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

func imagePageHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		// fetch image ID from url
		vars := mux.Vars(req)
		id := vars["id"]
		if len(id) > fs.ImgNameLength {
			if strings.Contains(id, ".") {
				id = strings.Split(id, ".")[0]
			}
			found, imgPath, imgID, title, errc, err := db.CheckIfImageExists(id)
			if err != nil {
				misc.PrintMessage(1, "http", "reqimg.go", "imagePageHandler()", string(errc)+err.Error())
				if errc == 500 {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				} else if errc == 404 {
					errorHandler(w, req)
					return
				}
			}
			if found == true {
				img := Img{title, "/img/" + imgID}
				fp := path.Join("public", "img.html")
				tmpl, err := template.ParseFiles(fp)
				if err != nil {
					misc.PrintMessage(1, "http", "reqimg.go", "imagePageHandler()", "500. Couldn't parse template.")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				// return (execute) the template or print an error if one occurs
				if err := tmpl.Execute(w, img); err != nil {
					misc.PrintMessage(1, "http", "reqimg.go", "imagePageHandler()", "500. Couldn't return template.")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				misc.PrintMessage(0, "http", "reqimg.go", "imagePageHandler()", imgPath+" > img.html")
				return
			}
			errorHandler(w, req)
			return
		}
		errorHandler(w, req)
		return
	}
}

func uploadHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		req.ParseMultipartForm(32 << 20)
		file, fileheader, err := req.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		var id string
		created := false
		for created != true {
			id = generateImageID()
			err = db.CheckImageID(id)
			if err != nil {
				if err.Error() == "Image ID '"+id+"' exists." {
					created = false
					continue
				} else {
					created = false
					misc.PrintMessage(1, "http", "reqimg.go", "uploadHandler()", err.Error())
					return
				}
			}
			created = true
		}

		err = imageMkdir(id, strings.Split(fileheader.Filename, ".")[1])
		if err != nil {
			misc.PrintMessage(1, "http", "reqimg.go", "uploadHandler()", err.Error())
			return
		}

		imgPath := id[0:fs.ImgStorageSubDirNameLength] + "/" +
			id[fs.ImgStorageSubDirNameLength:fs.ImgStorageSubDirNameLength*2] + "/" +
			id[fs.ImgStorageSubDirNameLength*2:fs.ImgStorageSubDirNameLength*3] + "/" +
			id[fs.ImgStorageSubDirNameLength*3:fs.ImgStorageSubDirNameLength*4] + "/" +
			id[((fs.ImgStorageSubDirNameLength*4)+1):len(id)] + "." + strings.Split(fileheader.Filename, ".")[1]

		err = db.StoreImage(id, strings.Split(fileheader.Filename, ".")[0], imgPath, strings.Split(fileheader.Filename, ".")[1])
		if err != nil {
			misc.PrintMessage(1, "http", "reqimg.go", "uploadHandler()", err.Error())
			return
		}

		http.Redirect(w, req, "/"+id, http.StatusFound)

		if id != "" {
			filePath := fs.ImgStoragePath +
				id[0:fs.ImgStorageSubDirNameLength] + "/" +
				id[fs.ImgStorageSubDirNameLength:fs.ImgStorageSubDirNameLength*2] + "/" +
				id[fs.ImgStorageSubDirNameLength*2:fs.ImgStorageSubDirNameLength*3] + "/" +
				id[fs.ImgStorageSubDirNameLength*3:fs.ImgStorageSubDirNameLength*4] + "/" +
				id[((fs.ImgStorageSubDirNameLength*4)+1):len(id)] + "." +
				strings.Split(fileheader.Filename, ".")[1]

			f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
			defer f.Close()
			if err != nil {
				misc.PrintMessage(1, "http", "reqimg.go", "uploadHandler()", err.Error())
				return
			}
			misc.PrintMessage(0, "http", "reqimg.go", "uploadHandler()", "File "+filePath+" has been created.")
			bytesCopied, err := io.Copy(f, file)
			if err != nil {
				misc.PrintMessage(1, "http", "reqimg.go", "uploadHandler()", err.Error())
				return
			}
			misc.PrintMessage(0, "http", "reqimg.go", "uploadHandler()", "Content of uploaded image ("+strconv.FormatInt(bytesCopied, 10)+" Bytes) has been copied to "+filePath+".")
		}
	}
}

func generateImageID() string {
	// generate a random hash
	// length is provided by an environment
	// variable or user input
	b := make([]byte, fs.ImgNameLength)
	rand.Read(b)
	filename := fmt.Sprintf("%x", b)
	return filename
}

func imageMkdir(id string, fileExt string) error {
	// check if the file already exists
	// to instantly generate a new name
	// before running through the whole process
	if _, err := os.Stat(
		fs.ImgStoragePath +
			id[0:fs.ImgStorageSubDirNameLength] + "/" +
			id[fs.ImgStorageSubDirNameLength:fs.ImgStorageSubDirNameLength*2] + "/" +
			id[fs.ImgStorageSubDirNameLength:fs.ImgStorageSubDirNameLength*3] + "/" +
			id[fs.ImgStorageSubDirNameLength*3:fs.ImgStorageSubDirNameLength*4] + "/" +
			id[((fs.ImgStorageSubDirNameLength*4)+1):len(id)] + "." + fileExt); os.IsNotExist(err) {
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

			// If the OFFSET is lower than 12
			// there are still directories to be created
			// since there need to be at least 4 directories
			if offset < fs.ImgStorageSubDirNameLength*4 {

				// Check if the directory that is to be created
				// already exists.
				//
				// If it doesn't the directory will be created,
				// the directory will be added to the current path
				// and the offset will be increased by the number of
				// letters used (3) from the filename random hash
				if _, err := os.Stat(fs.ImgStoragePath + currentPath + id[offset:offset+fs.ImgStorageSubDirNameLength] + "/"); os.IsNotExist(err) {
					// doesn't exist
					misc.PrintMessage(2, "http", "reqimg.go", "imageMkdir()", fs.ImgStoragePath+currentPath+id[offset:offset+fs.ImgStorageSubDirNameLength]+"/ created.")
					os.Mkdir(fs.ImgStoragePath+currentPath+id[offset:offset+fs.ImgStorageSubDirNameLength]+"/", 0776)
					currentPath += id[offset:offset+fs.ImgStorageSubDirNameLength] + "/"
					offset += fs.ImgStorageSubDirNameLength
				} else {
					// exists
					misc.PrintMessage(2, "http", "reqimg.go", "imageMkdir()", fs.ImgStoragePath+currentPath+id[offset:offset+fs.ImgStorageSubDirNameLength]+"/ already existed and thus is not created!")
					currentPath += id[offset:offset+fs.ImgStorageSubDirNameLength] + "/"
					offset += fs.ImgStorageSubDirNameLength
				}

				// If the OFFSET is higher than 9
				// the next step is to create the file itself.
				// The function won't do this itself.
				// Instead it will check, if the file exists (again, just to be sure)
				// and if the file doesn't exist, it will return the name of the file
				// and an empty error.
			} else {
				if _, err := os.Stat(fs.ImgStoragePath + currentPath + id[offset:len(id)] + "." + fileExt); os.IsNotExist(err) {
					// doesn't exist
					return nil
				}
			}
		}
	}
	return errors.New("Something went wrong. (findAvailableName)")
}

func favIcoHandler(w http.ResponseWriter, req *http.Request) {
	fp := path.Join("public/img/", "favicon.ico")
	http.ServeFile(w, req, fp)
}
