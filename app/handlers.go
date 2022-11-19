package main

import (
	"fmt"
	"net/http"
	"strconv"

	"blogdemo.batou.cn/common/validator"
	"github.com/go-chi/chi/v5"
)

type blogForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	validator.Validator `form:"_"`
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

	var form blogForm
	app.formDecoder.Decode(&form, r.PostForm)

	form.CheckField(validator.NotBlank(form.Title), "title", "标题不能为空")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "标题不能超过100字")
	form.CheckField(validator.NotBlank(form.Content), "content", "内容不能为空")

	if !form.Valid() {
		data := app.newTemplateData()
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "blog_create.html", data)
		return
	}

	id, err := app.blogs.Insert(form.Title, form.Content)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/blog/view/%d", id), http.StatusSeeOther)
}
