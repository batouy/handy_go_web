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

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.checkAuth)

	route.Get("/", handlerToFunc(dynamic.ThenFunc(app.home)))
	route.Get("/about", handlerToFunc(dynamic.ThenFunc(app.about)))
	route.Get("/blog/view/{id}", handlerToFunc(dynamic.ThenFunc(app.blogView)))

	route.Get("/login", handlerToFunc(dynamic.ThenFunc(app.login)))
	route.Post("/login", handlerToFunc(dynamic.ThenFunc(app.loginPost)))
	route.Get("/register", handlerToFunc(dynamic.ThenFunc(app.register)))
	route.Post("/register", handlerToFunc(dynamic.ThenFunc(app.registerPost)))

	auth := dynamic.Append(app.auth)
	route.Get("/blog/create", handlerToFunc(auth.ThenFunc(app.blogCreate)))
	route.Post("/blog/store", handlerToFunc(auth.ThenFunc(app.blogStore)))
	route.Post("/logout", handlerToFunc(auth.ThenFunc(app.logout)))

	commonMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return commonMiddleware.Then(route)
}

func handlerToFunc(handler http.Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}
