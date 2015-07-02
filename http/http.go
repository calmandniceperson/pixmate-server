/*
 * IMGCAT
 * GO PROTOTYPE
 * 2015
 */

package http

import (
  "net/http"
  "github.com/gorilla/mux"
  "github.com/fatih/color"
  "github.com/codegangsta/negroni"
)

/*
 * initialise http server
 */
func Init(){
  r := mux.NewRouter().StrictSlash(false)

  /*
   * main page
   */
  r.HandleFunc("/", mainPageHandler)

  /*
   * custom profile url
   */
  r.HandleFunc("/people/{name}", peoplePageHandler)

  /*
   * init negroni middleware
   */
  n := negroni.New(
    negroni.NewRecovery(),
    negroni.HandlerFunc(MiddleWare),
    //negroni.NewLogger(),
    negroni.NewStatic(http.Dir("public")),
  )
  n.UseHandler(r)
  color.Green("IMGCAT Server running on port 8000")
  n.Run(":8000")
}
