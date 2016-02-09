package http

import (
	"imgturtle/db"
	"imgturtle/fsys"
	"imgturtle/io"
	"net/http"
	"strconv"
	"strings"
)

func apiUploadHandler(w http.ResponseWriter, req *http.Request) {
	cio.PrintMessage(0, "Request from pixmate client...")
	if req.Method == "POST" {
		req.ParseMultipartForm(0)
		file, fileheader, err := req.FormFile("image")
		title := req.FormValue("title")
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
		err = db.StoreImage(id, title, filePath, strings.Split(fileheader.Filename, ".")[1] /*extension*/)
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
		w.Write([]byte(id))
	}

}
