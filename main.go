package main

import (
	"io"
	"net/http"
  "fmt"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
  fmt.Println("User called for /")
}

func main() {
  fmt.Println("Server running on port 8000")
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8000", nil)
}
