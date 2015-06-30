package reqh

import (
  "fmt"
  "net/http"
)

func Init(){
  fmt.Println("Server running on port 8000")
	http.HandleFunc("/", mainPageHandler)
	http.ListenAndServe(":8000", nil)
}

func mainPageHandler(w http.ResponseWriter, req *http.Request) {
	resourcePath := "public/html" + req.URL.Path + ".html"
	if req.URL.Path == "/" {
		resourcePath = "public/html/" + "welcome.html"
	}
	fmt.Printf("serving static file => %s", resourcePath)
	http.ServeFile(w, req, resourcePath)
}
