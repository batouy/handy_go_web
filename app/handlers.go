package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	blogs, err := app.blogs.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData()
	data.Blogs = blogs

	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) blogView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.clientError(w, http.StatusNotFound)
	}

	blog, err := app.blogs.Get(id)
	if err != nil {
		app.clientError(w, http.StatusNotFound)
		return
	}

	data := app.newTemplateData()
	data.Blog = blog

	app.render(w, http.StatusOK, "blog.html", data)
}

func (app *application) blogCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	app.render(w, http.StatusOK, "blog_create.html", data)
}

func (app *application) blogStore(w http.ResponseWriter, r *http.Request) {
	id, err := app.blogs.Insert("测试4", "测试内容44444")

	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/blog/view/%d", id), http.StatusSeeOther)
}
