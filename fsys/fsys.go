package fsys

import (
	"bufio"
	"os"
	"pixmate-server/io"
	"strconv"

	"github.com/fatih/color"
)

// ImgStorage
// Stores the image storage's path in the file system
// to make it accessible from other packages
var ImgStoragePath string

// ImgStorageSubDirNameLength
// Stores the length of the subdirectory names
// in the image storage
var ImgStorageSubDirNameLength int

// ImgNameLength
// Stores the lengths of the images' names
var ImgNameLength int

func Start() {
	if os.Getenv("IMGSTORAGE_LOCATION") != "" {
		ImgStoragePath = os.Getenv("IMGSTORAGE_LOCATION")
	} else {
		reader := bufio.NewReader(os.Stdin)
		color.Cyan("Enter location of image storage: ")
		ImgStoragePath, _ = reader.ReadString('\n')
	}
	cio.PrintMessage(1, "Image storage is being created...")
	if _, err := os.Stat(ImgStoragePath); os.IsNotExist(err) {
		// doesn't exist
		os.Mkdir(ImgStoragePath, 0776)
		cio.PrintMessage(2, (ImgStoragePath + " created."))
	} else {
		cio.PrintMessage(2, (ImgStoragePath + " already existed."))
	}

	if os.Getenv("IMG_NAME_LENGTH") != "" {
		ImgNameLength, _ = strconv.Atoi(os.Getenv("IMG_NAME_LENGTH"))
	} else {
		reader := bufio.NewReader(os.Stdin)
		color.Cyan("Enter length of image names: ")
		temp, _ := reader.ReadString('\n')
		ImgNameLength, _ = strconv.Atoi(temp)
	}
}
