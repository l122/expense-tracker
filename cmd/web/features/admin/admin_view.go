package admin

import (
	"html/template"
	"net/http"
)

type AdminView struct {
	templ *template.Template
}

func NewAdminView(templ *template.Template) *AdminView {
	return &AdminView{
		templ: templ,
	}
}

func (i *AdminView) Index(w http.ResponseWriter) {
	if err := i.templ.ExecuteTemplate(w, "admin", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
