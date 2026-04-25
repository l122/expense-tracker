package web

import (
	"html/template"
	"net/http"
)

type IndexView struct {
	templ *template.Template
}

func NewIndexView(templ *template.Template) *IndexView {
	return &IndexView{templ: templ}
}

func (i *IndexView) Index(w http.ResponseWriter) {
	if err := i.templ.ExecuteTemplate(w, "index", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
