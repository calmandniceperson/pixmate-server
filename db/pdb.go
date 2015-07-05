package db

import (
  "log"
  "github.com/fatih/color"
  "database/sql"
	_ "github.com/lib/pq"
)

func Init(){
  db, err := sql.Open("postgres", "user=imgturtlea password=imgtadmin dbname=imgturtle sslmode=disable")
  if err != nil {
    color.Red("ERR: pdb.go Init() => PostgreSQL config could not be established.")
    color.Red(err.Error())
  }
  defer db.Close()

  /*
  // test connection
  err = db.Ping()
  if err != nil{ // connection not successful
    color.Red("ERR: pdb.go Init() => Database connection not working.")
    color.Red(err.Error())
  }else{ // connection successful
    /*
     * variables for storing the user data
     *
    var (
      id int
      uname string
      pw string
    )

    /*
     * run query
     *
    rows, err := db.Query(`select user_id, user_name, user_pw from "ITUser"`)
    if err != nil{
      color.Red("ERR: pdb.go Init() => Query could not be executed.")
      color.Red(err.Error())
    }

    /*
     * close database connection at the
     * end of the enclosing function
     *
    defer rows.Close()

    /*
     * .Next() prepares the next data column for reading
     * .Scan(values) transfers the data to the given variables
     *
    for rows.Next() {
      err := rows.Scan(&id, &uname, &pw)
      if err != nil {
        color.Red("ERR: pdb.go Init() => Fetched values could not be scanned.")
        color.Red(err.Error())
      }
      log.Println(id, uname, pw)
    }

    err = rows.Err()
    if err != nil {
      color.Red("ERR: pdb.go Init() => An error occured.")
      color.Red(err.Error())
    }
    */
  }
}
