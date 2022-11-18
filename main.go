package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("homepage"))
}

func blogView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 1 {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte(fmt.Sprintf("view blog id: %d", id)))
}

func blogCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
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
