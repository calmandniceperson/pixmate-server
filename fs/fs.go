package fs

import (
	"bufio"
	"os"

	"github.com/fatih/color"
)

// ImgStorage
// Stores the image storage's path in the file system
// to make it accessible from other packages
var ImgStoragePath string

// CreateImageStorageIfNotExists will check whether there is already
// an image storage path in the environment variables.
// If there is none, the user of this server will be asked
// to provide a path.
// This image storage path will be used in http/reqh.go
// to store images.
func CreateImageStorageIfNotExists() {
	if os.Getenv("IMGSTORAGE_LOCATION") != "" {
		ImgStoragePath = os.Getenv("IMGSTORAGE_LOCATION")
	} else {
		reader := bufio.NewReader(os.Stdin)
		color.Cyan("Enter location of image storage: ")
		ImgStoragePath, _ = reader.ReadString('\n')
	}

	if _, err := os.Stat(ImgStoragePath); os.IsNotExist(err) {
		// doesn't exist
		os.Mkdir(ImgStoragePath, 0776)
		color.Cyan("INF: " + ImgStoragePath + " created.")
	} else {
		color.Cyan("INF: " + ImgStoragePath + " aleady existed.")
	}
}
