package main

import (
	"log"
	"net/http"
)

func main() {
	// serve files from the "public" directory
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	log.Println("Serving generated static site on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
