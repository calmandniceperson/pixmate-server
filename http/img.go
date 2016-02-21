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
	"time"

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
			found, imgPath, imgID, title, tUp, ttlTime, ttlViews, errc, err := db.GetImage(id)

			/* Check whether the time to live limit
			 * has been reached (time of upload + allowed minutes)
			 */
			if ttlTime > 0 {
				ttlTimeT := time.Date(
					tUp.Year(),
					tUp.Month(),
					tUp.Day(),
					tUp.Hour(),
					tUp.Minute()+int(ttlTime),
					tUp.Second(),
					tUp.Nanosecond(),
					time.Local,
				)
				if time.Now().After(ttlTimeT) /*|| ttlViews == 0*/ {
					http.Redirect(w, req, "/error", 403)
					err := db.DeleteImage(imgID)
					if err != nil {
						cio.PrintMessage(1, err.Error())
					}
					return
				}
			}

			/* Decrease the views left for this image
			 * if the view count is not unlimited (unlimited = smaller than 0)
			 * Check whether the view limit has been
			 * reached.
			 */
			print(ttlViews)
			if ttlViews > -1 {
				if ttlViews == 0 {
					http.Redirect(w, req, "/error", 403)
					err := db.DeleteImage(imgID)
					if err != nil {
						cio.PrintMessage(1, err.Error())
					}
					return
				}
				err := db.UpdateImageViewCount(imgID)
				if err != nil {
					cio.PrintMessage(1, err.Error())
				}
				ttlViews--
			}

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
			found, imgPath, imgID, _, tUp, ttlTime, _, errc, err := db.GetImage(id)

			/* Check whether the time to live limit
			 * has been reached (time of upload + allowed minutes)
			 */
			if ttlTime > 0 {
				ttlTimeT := time.Date(
					tUp.Year(),
					tUp.Month(),
					tUp.Day(),
					tUp.Hour(),
					tUp.Minute()+int(ttlTime),
					tUp.Second(),
					tUp.Nanosecond(),
					time.Local,
				)
				if time.Now().After(ttlTimeT) /*|| ttlViews == 0*/ {
					http.Redirect(w, req, "/error", 403)
					err := db.DeleteImage(imgID)
					if err != nil {
						cio.PrintMessage(1, err.Error())
					}
					return
				}
			}

			/* There is no check for the view count
				 * on direct calls for the image itself (instead
			   * of its page) to avoid the count decreasing
				 * when someone downloads the image
			*/

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

		/* Retrieve the file */
		file, fileheader, err := req.FormFile("uploadFile")
		if err != nil {
			cio.PrintMessage(1, err.Error())
			return
		}

		/* Retrieve:
		*		- file extension
		*		- title: if there is a custom title, it will be used, otherwise the
								   current file name will be used
		*		- time to live (per time), if set
		*		- time to live (per views), if set
		*/
		fileExt := strings.Split(fileheader.Filename, ".")[1]
		var title string
		if req.FormValue("title") != "" {
			title = req.FormValue("title")
		} else {
			title = strings.Split(fileheader.Filename, ".")[0]
		}
		var ttlTime int64
		var ttlViews int64
		if req.FormValue("ttlTime") != "" {
			ttlTime, err = strconv.ParseInt(req.FormValue("ttlTime"), 10, 64)
			if err != nil {
				cio.PrintMessage(1, err.Error())
			}
		} else {
			ttlTime = 60 * 24 * 7 * 4 * 2 /*2 months (default)*/
		}
		if req.FormValue("ttlViews") != "" {
			ttlViews, err = strconv.ParseInt(req.FormValue("ttlViews"), 10, 64)
			if ttlViews == 0 {
				// do not allow it to be 0 from the very start
				// make it unlimited instead
				ttlViews = -1
			}
			if err != nil {
				cio.PrintMessage(1, err.Error())
			}
		} else {
			ttlViews = -1
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

		/* Generate the path where the image will be stored
		 * from the image storage directory path + the image id (generated by generateImageID()) +
		 * file extension
		 */
		filePath := fsys.ImgStoragePath + id + "." + fileExt

		/* Store the image
		 * Options:
		 *		- custom title
		 *		- custom time to live (per time)
		 *		- custom time to live (per views)
		 */
		err = db.StoreImage(id, title, filePath, fileExt, ttlTime, ttlViews)
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
