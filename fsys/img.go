package fsys

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"pixmate-server/db"
	"pixmate-server/io"
	"strings"
	"time"
)

func StoreImage(filePath string, file multipart.File) (int64, error) {
	newFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	defer newFile.Close()
	bytesCopied, err := io.Copy(newFile, file)
	if err != nil {
		cio.PrintMessage(1, err.Error())
		return 0, err
	}
	return bytesCopied, nil
}

func RemoveOldImages() error {
	files, err := ioutil.ReadDir(os.Getenv("IMGSTORAGE_LOCATION"))
	if err != nil {
		cio.PrintMessage(1, err.Error())
	}
	for _ /*index*/, ele := range files {
		isOld, err := db.CheckImgTTLExceeded(strings.Split(ele.Name(), ".")[0])
		if err != nil {
			cio.PrintMessage(1, err.Error())
			return err
		}
		if isOld {
			err := db.DeleteImage(strings.Split(ele.Name(), ".")[0])
			if err != nil {
				cio.PrintMessage(1, err.Error())
				return err
			}
			err = DeleteFile(ele.Name())
			if err != nil {
				cio.PrintMessage(1, err.Error())
			}
		}
	}
	time.Sleep((60 * 2) * time.Minute)
	RemoveOldImages()
	return nil
}

func DeleteFile(filename string) error {
	err := os.Remove(os.Getenv("IMGSTORAGE_LOCATION") + filename)
	if err != nil {
		return err
	}
	cio.PrintMessage(2, "Deleted file "+filename)
	return nil
}
