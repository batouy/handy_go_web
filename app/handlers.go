package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	blogs, err := app.blogs.Latest()
	if err != nil {
		app.errorLog.Print(err.Error())
		http.Error(w, "服务内部出错", http.StatusInternalServerError)
		return
	}

	for _, blog := range blogs {
		fmt.Fprintf(w, "%+v\n", blog)
	}

	return

	files := []string{
		"./resources/views/layouts/default.html",
		"./resources/views/partials/nav.html",
		"./resources/views/home.html",
	}

	tp, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Print(err.Error())
		http.Error(w, "服务内部出错", http.StatusInternalServerError)
		return
	}

	err = tp.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		app.errorLog.Print(err.Error())
		http.Error(w, "服务内部出错", http.StatusInternalServerError)
	}
}

func (app *application) blogView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 1 {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte(fmt.Sprintf("view blog id: %d", id)))
}

func (app *application) blogCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("blog create"))
}
