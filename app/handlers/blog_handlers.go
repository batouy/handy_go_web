package handlers

import (
	"fmt"
	"net/http"
	"strconv"
)

func BlogView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 1 {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte(fmt.Sprintf("view blog id: %d", id)))
}

func BlogCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("blog create"))
}
