package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/book/view", app.bookView)
	mux.HandleFunc("/book/create", app.bookCreate)
	mux.HandleFunc("/book/delete", app.bookDelete)
	mux.HandleFunc("/book/edit", app.bookEdit)

	return app.secureHeaders(mux)
}
