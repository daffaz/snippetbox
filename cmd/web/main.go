package main

import (
	"log"
	"net/http"
)

func main() {
	// initialize serve mux
	mux := http.NewServeMux()

	// for static file
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// root
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("Running in port :4000")
	log.Fatal(http.ListenAndServe(":4000", mux))
}
