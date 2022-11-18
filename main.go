package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("homepage"))
}

func blogView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("blog view"))
}

func blogCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	w.Write([]byte("blog create"))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/blog/view", blogView)
	mux.HandleFunc("/blog/create", blogCreate)

	log.Print("启动web服务，端口4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
