package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	// initialize serve mux
	mux := http.NewServeMux()

	// for static file
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// root
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// static file route
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
