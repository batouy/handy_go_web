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

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	route.Get("/", handlerToFunc(dynamic.ThenFunc(app.home)))
	route.Get("/blog/view/{id}", handlerToFunc(dynamic.ThenFunc(app.blogView)))
	route.Get("/blog/create", handlerToFunc(dynamic.ThenFunc(app.blogCreate)))
	route.Post("/blog/store", handlerToFunc(dynamic.ThenFunc(app.blogStore)))

	commonMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return commonMiddleware.Then(route)
}

func handlerToFunc(handler http.Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}
