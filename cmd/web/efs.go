package web

import (
	"embed"
	"html/template"
	"log"
)

//go:embed features/*/*.gohtml features/*/*/*.gohtml
var Files embed.FS

func ParseTemplates() *template.Template {
	// Add patterns for the root directory and feature subdirectories
	templates, err := template.ParseFS(Files, "features/*/*.gohtml", "features/*/*/*.gohtml")
	if err != nil {
		log.Fatal(err)
	}
	return templates
}
