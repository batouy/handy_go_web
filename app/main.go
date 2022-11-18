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

	tp, err := template.ParseFiles("./resources/views/home.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "服务内部出错", http.StatusInternalServerError)
		return
	}

	err = tp.Execute(w, nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "服务内部出错", http.StatusInternalServerError)
		return
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/blog/view", handlers.BlogView)
	mux.HandleFunc("/blog/create", handlers.BlogCreate)

	log.Print("启动web服务，端口4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
