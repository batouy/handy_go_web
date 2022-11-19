package main

import (
	"path/filepath"
	"text/template"
	"time"
)

func humanDate(t time.Time) string {
	return t.Format("2006:01:02 15:04:05")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

// 在内存中缓存页面
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./resources/views/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		files := []string{
			"./resources/views/layouts/default.html",
			"./resources/views/partials/nav.html",
			page,
		}

		tp, err := template.New(name).Funcs(functions).ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cache[name] = tp
	}

	return cache, nil
}
