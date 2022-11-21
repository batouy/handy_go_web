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

type regForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"_"`
}

type loginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
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

	data := app.newTemplateData(r)
	data.Blogs = blogs

	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "about.html", data)
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

	data := app.newTemplateData(r)
	data.Blog = blog

	app.render(w, http.StatusOK, "blog.html", data)
}

func (app *application) blogCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
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
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "blog_create.html", data)
		return
	}

	id, err := app.blogs.Insert(form.Title, form.Content)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "恭喜，文章保存成功！")

	http.Redirect(w, r, fmt.Sprintf("/blog/view/%d", id), http.StatusSeeOther)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = loginForm{}
	app.render(w, http.StatusOK, "login.html", data)
}

func (app *application) loginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var form loginForm
	err = app.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "邮箱不能为空")
	form.CheckField(validator.IsEmail(form.Email, validator.EmailRX), "email", "邮箱不能为空")
	form.CheckField(validator.NotBlank(form.Password), "password", "密码不能为空")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)

	if err != nil {
		form.AddNonFieldError("账号或密码有误")
	}

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	app.sessionManager.RenewToken(r.Context())
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	app.sessionManager.Put(r.Context(), "flash", "登录成功！")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) register(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = regForm{}
	app.render(w, http.StatusOK, "register.html", data)
}

func (app *application) registerPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var form regForm
	app.formDecoder.Decode(&form, r.PostForm)

	form.CheckField(validator.NotBlank(form.Name), "name", "用户名不能为空")
	form.CheckField(validator.NotBlank(form.Email), "email", "邮箱不能为空")
	form.CheckField(validator.IsEmail(form.Email, validator.EmailRX), "email", "邮箱非法")
	form.CheckField(validator.NotBlank(form.Password), "password", "密码不能为空")
	form.CheckField(validator.MinChars(form.Password, 6), "password", "密码至少6个字符")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "register.html", data)
		return
	}

	_, err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	app.sessionManager.RenewToken(r.Context())
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "退出成功！")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
