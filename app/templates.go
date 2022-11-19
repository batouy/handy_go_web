package main

import (
	"text/template"
	"time"
)

func humanDate(t time.Time) string {
	return t.Format("2006:01:02 15:04:05")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}
