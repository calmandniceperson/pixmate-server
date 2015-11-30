package fsys

import (
	"imgturtle/io"
	"io"
	"mime/multipart"
	"os"
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
