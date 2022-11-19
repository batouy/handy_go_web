package main

import (
	"fmt"
	"net/http"
	"strconv"
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

	data := app.newTemplateData()
	data.Blogs = blogs

	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) blogView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 1 {
		http.NotFound(w, r)
		return
	}

	blog, err := app.blogs.Get(id)
	if err != nil {
		app.errorLog.Print(err.Error())
		http.Error(w, "未查到数据", http.StatusNotFound)
		return
	}

	data := app.newTemplateData()
	data.Blog = blog

	app.render(w, http.StatusOK, "blog.html", data)
}

func (app *application) blogCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := app.blogs.Insert("测试4", "测试内容44444")

	if err != nil {
		app.errorLog.Print(err.Error())
		http.Error(w, "内部出错", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/blog/view?id=%d", id), http.StatusSeeOther)
}
