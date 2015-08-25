package db

import (
	"bufio"
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/pbkdf2"
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
		color.Red("ERR: pdb.go Init() => PostgreSQL config could not be established.")
		color.Red(err.Error())
	}

	// test connection
	err = db.Ping()
	if err != nil { // connection not successful
		color.Red("ERR: pdb.go Init() => Database connection not working.")
		log.Fatal(err.Error())
	}
}

// CheckUserCredentials handles the database part of the
// login process
func CheckUserCredentials(ue string, pwd string) (bool, error) {
	rows, err := db.Query("select user_name, user_email, user_pw, user_hash from imgturtle.user where user_name='" + ue + "' or user_email='" + ue + "'")
	if err != nil {
		color.Red("ERR@pdb.go@CheckUserCredentials() => %s", err.Error())
		return false, err
	}

	if rows != nil {
		defer rows.Close()

		var (
			fUname string
			fEmail string
			fPw    string
			fHash  string
		)
		if rows.Next() {
			err := rows.Scan(&fUname, &fEmail, &fPw, &fHash)
			if err != nil {
				color.Red("ERR: pdb.go CheckUserCredentials => Fetched values could not be scanned.")
				color.Red(err.Error())
				return false, err
			}
			if fPw == string(pbkdf2.Key([]byte(pwd), []byte(fHash), 4096, 32, sha1.New)) {
				color.Green("User %s entered a valid password.", fUname)
				return true, nil
			}
			color.Red("User %s entered an invalid password.", fUname)
			return false, errors.New("Incorrect password.")
		}
		color.Green("User %s could not be found.", ue)
		return false, errors.New("No such user.")
	}
	return false, nil
}

// InsertNewUser handles the database part of the process of
// registering a new user
func InsertNewUser(uname string, pwd string, email string) error {
	rows, err := db.Query("select user_name, user_email from imgturtle.user where user_name='" + uname + "' or user_email='" + email + "'")
	if err != nil {
		color.Red("ERR@pdb.go@InsertNewUser() => %s", err.Error())
	}

	if rows != nil {
		defer rows.Close()

		var (
			funame string
			femail string
		)
		for rows.Next() {
			err := rows.Scan(&funame, &femail)
			if err != nil {
				color.Red("ERR: pdb.go InsertNewUser() => Fetched values could not be scanned.")
				color.Red(err.Error())
				return err
			}
			if funame == uname && femail == email {
				return errors.New("User name '" + uname + "' and e-mail address '" + email + "' in use.")
			} else if funame == uname {
				return errors.New("User name '" + uname + "' in use.")
			} else if femail == email {
				return errors.New("E-mail address '" + email + "'in use.")
			}
		}
	}

	b := make([]byte, 32)
	rand.Read(b)
	salt := fmt.Sprintf("%x", b)

	epw := pbkdf2.Key([]byte(pwd), []byte(salt), 4096, 32, sha1.New)

	stmt, err := db.Prepare("INSERT INTO imgturtle.user(user_name,user_pw,user_email,user_hash) VALUES($1,$2,$3,$4)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(uname, epw, email, salt)
	if err != nil {
		return err
	}

	return nil
}

func CheckIfImageExists(id string) (bool, string, string, error) {
	rows, err := db.Query("select image_id, image_title, image_f_ext from imgturtle.img where image_id='" + id + "'")
	if err != nil {
		color.Red("ERR@pdb.go@InsertNewUser() => %s", err.Error())
	}
	var (
		fid  string
		ftit string
		fext string
	)

	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&fid, &ftit, &fext)
			if err != nil {
				color.Red("ERR: pdb.go CheckIfImageExists() => Fetched values could not be scanned.")
				color.Red(err.Error())
				return false, "", "", err
			}
		}
		if fid == id {
			return true, ftit, fext, nil
		}
	}
	return false, "", "", errors.New("No image with id " + id + " could be found.")
}

func CheckImageID(id string) error {
	rows, err := db.Query("select image_id from imgturtle.img where image_id='" + id + "'")
	if err != nil {
		color.Red("ERR@pdb.go@InsertNewUser() => %s", err.Error())
	}

	if rows != nil {
		defer rows.Close()

		var fid string
		for rows.Next() {
			err := rows.Scan(&fid)
			if err != nil {
				color.Red("ERR: pdb.go InsertNewUser() => Fetched values could not be scanned.")
				color.Red(err.Error())
				return err
			}
			if fid == id {
				return errors.New("Image ID '" + id + "' in use.")
			}
		}
	}
	return nil
}

// StoreImage stores all of an image's information in the database
func StoreImage(id string, title string, ext string /*, desc string, uploader_id string, uploader_name string*/) error {
	stmt, err := db.Prepare("INSERT INTO imgturtle.Img(image_id, image_title, image_f_ext) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id, title, ext)
	if err != nil {
		return err
	}

	return nil
}
