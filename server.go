package main

import (
	"net/http"

	"github.com/gobuffalo/packr/v2"
)

// Serve index.html on "/"
// Serve all other assets on

func main() {
	box := packr.New("ReactApp", "./web/build")

	http.Handle("/", http.FileServer(box))
	http.ListenAndServe(":8080", nil)
}
