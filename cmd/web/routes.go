package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// middleware chain
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// dynamic middleware
	dynamiMiddleware := alice.New(app.session.Enable)

	// initialize serve mux
	mux := pat.New()

	// for static file
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// route
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snippet/create", dynamiMiddleware.ThenFunc(http.HandlerFunc(app.createSnippetForm)))
	mux.Post("/snippet/create", dynamiMiddleware.ThenFunc(http.HandlerFunc(app.createSnippet)))
	mux.Get("/snippet/:id", dynamiMiddleware.ThenFunc(http.HandlerFunc(app.showSnippet)))
	// static file route
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
