package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func (app *application) render(w http.ResponseWriter, status int, page string, data any) {
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
