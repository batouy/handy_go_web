package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	route := chi.NewRouter()

	route.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		fileServer := http.FileServer(http.Dir("./resources/asserts/"))
		http.StripPrefix("/static", fileServer).ServeHTTP(w, r)
	})

	route.Get("/", app.home)
	route.Get("/blog/view/{id}", app.blogView)
	route.Get("/blog/create", app.blogCreate)
	route.Post("/blog/store", app.blogStore)

	commonMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return commonMiddleware.Then(route)
}
