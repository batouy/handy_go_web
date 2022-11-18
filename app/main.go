package main

import (
	"log"
	"net/http"
	"text/template"

	"blogdemo.batou.cn/app/handlers"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./resources/views/layouts/default.html",
		"./resources/views/partials/nav.html",
		"./resources/views/home.html",
	}

	tp, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "服务内部出错", http.StatusInternalServerError)
		return
	}

	err = tp.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "服务内部出错", http.StatusInternalServerError)
		return
	}
}

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./resources/asserts/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/blog/view", handlers.BlogView)
	mux.HandleFunc("/blog/create", handlers.BlogCreate)

	log.Print("启动web服务，端口4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
