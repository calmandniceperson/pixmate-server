package db

import (
	"bufio"
	"database/sql"
	"errors"
	"imgturtle/io"
	"os"

	"github.com/fatih/color"
	_ "github.com/lib/pq"
)

var db *sql.DB

// Start is the database package launch method
// it enters or fetches the data required for the database
func Start() {
	/*
	 * allow user to enter db data
	 * used instead of environment variables
	 * if there are none
	 * since the service is open source
	 */
	var (
		uname string
		pw    string
		name  string
	)
	if os.Getenv("DB_UNAME") == "" && os.Getenv("DB_NAME") == "" {
		reader := bufio.NewReader(os.Stdin)
		color.Cyan("Enter db user name: ")
		uname, _ = reader.ReadString('\n')

		color.Cyan("Enter db pw: ")
		pw, _ = reader.ReadString('\n')

		color.Cyan("Enter db name: ")
		name, _ = reader.ReadString('\n')
	} else {
		uname = os.Getenv("DB_UNAME")
		pw = os.Getenv("DB_PW")
		name = os.Getenv("DB_NAME")
	}
	var err error
	db, err = sql.Open("postgres",
		"user="+uname+
			" password="+pw+
			" dbname="+name+
			" sslmode=disable")

	if err != nil {
		cio.PrintMessage(1, err.Error())
		return
	}
	// test connection
	err = db.Ping()
	if err != nil { // connection not successful
		cio.PrintMessage(1, err.Error())
		var rundb string
		reader := bufio.NewReader(os.Stdin)
		color.Cyan("Do you want to run the server without a working database module? (y/n) ")
		rundb, _ = reader.ReadString('\n')
		if rundb != "y\n" && rundb != "Y\n" {
			os.Exit(-1)
		}
	}
}

func CheckIfImageExists(id string) (bool, string, string, string, int, error) {
	rows, err := db.Query("SELECT image_id, image_title, image_path FROM imgturtle.img WHERE image_id='" + id + "'")
	if err != nil {
		cio.PrintMessage(1, err.Error())
	}
	var (
		fID   string
		fTit  string
		fPath string
	)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&fID, &fTit, &fPath)
			if err != nil {
				cio.PrintMessage(1, err.Error())
				return false, "", "", "", 500, err
			}
		}
		if fID == id {
			return true, fPath, fID, fTit, 200, nil
		}
	}
	return false, "", "", "", 404, errors.New("No image with ID " + id + " could be found.")
}

func CheckIfImageIDInUse(id string) error {
	rows, err := db.Query("SELECT image_id FROM imgturtle.img WHERE image_id='" + id + "'")
	if err != nil {
		cio.PrintMessage(1, err.Error())
	}
	if rows != nil {
		defer rows.Close()
		var fID string
		for rows.Next() {
			err := rows.Scan(&fID)
			if err != nil {
				cio.PrintMessage(1, err.Error())
				return err
			}
			if fID == id {
				return errors.New("Image ID '" + id + "' in use.")
			}
		}
	}
	return nil
}

// StoreImage stores all of an image's information in the database
func StoreImage(id string, title string, imgPath string, ext string /*, desc string, uploader_id string, uploader_name string*/) error {
	stmt, err := db.Prepare("INSERT INTO imgturtle.Img(image_id, image_title, image_path, image_f_ext) VALUES($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id, title, imgPath, ext)
	if err != nil {
		return err
	}
	return nil
}
