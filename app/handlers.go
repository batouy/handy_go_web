package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
)

type blogForm struct {
	Title       string
	Content     string
	FieldErrors map[string]string
}

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
	data.Form = blogForm{}
	app.render(w, http.StatusOK, "blog_create.html", data)
}

func (app *application) blogStore(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	fieldErrors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) > 100 {
		fieldErrors["title"] = "标题不能超过100字"
	}

	if strings.TrimSpace(content) == "" {
		fieldErrors["content"] = "内容不能为空"
	}

	form := blogForm{
		Title:       title,
		Content:     content,
		FieldErrors: fieldErrors,
	}

	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData()
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "blog_create.html", data)
		return
	}

	id, err := app.blogs.Insert(title, content)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/blog/view/%d", id), http.StatusSeeOther)
}
