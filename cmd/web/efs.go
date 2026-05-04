package web

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"strings"
)

//go:embed styles.css features components
var Files embed.FS

func ParseTemplates() *template.Template {
	var filenames []string

	// Recursively walk the embedded filesystem to find all .gohtml files
	err := fs.WalkDir(Files, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".gohtml") {
			filenames = append(filenames, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Parse all discovered template files
	templates, err := template.ParseFS(Files, filenames...)
	if err != nil {
		log.Fatal(err)
	}
	return templates
}
