package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// 执行页面渲染
func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	tp, ok := app.templateCache[page]
	if !ok {
		app.errorLog.Printf(fmt.Sprintf("the template %s does not exist", page))
		w.Write([]byte("页面渲染出错"))
		return
	}

	// 服务端先将页面渲染到 buffer，如果出错，提前捕获，不会将残缺信息展示到用户端
	buf := new(bytes.Buffer)

	err := tp.ExecuteTemplate(buf, "layout", data)
	if err != nil {
		app.errorLog.Printf(fmt.Sprintf("the template %s does not exist", page))
		w.Write([]byte("页面渲染出错"))
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:    time.Now().Year(),
		Flash:          app.sessionManager.PopString(r.Context(), "flash"),
		IsAutheticated: app.isAuthenticated(r),
		CSRFToken:      nosurf.Token(r),
	}
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuth, ok := r.Context().Value(isAuthenticated).(bool)
	if !ok {
		return false
	}
	return isAuth
}
