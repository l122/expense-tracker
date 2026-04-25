package web

import (
	"embed"
	"html/template"
	"log"
)

//go:embed *.gohtml
var Files embed.FS

func ParseTemplates() *template.Template {
	templates, err := template.ParseFS(Files, "*.gohtml")
	if err != nil {
		log.Fatal(err)
	}
	return templates
}
