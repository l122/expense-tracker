package navbar

import (
	"html/template"
	"net/http"
)

type View struct {
	templ *template.Template
}

func NewView(templ *template.Template) *View {
	return &View{
		templ: templ,
	}
}

func (i *View) Index(w http.ResponseWriter) {
	if err := i.templ.ExecuteTemplate(w, "navbar", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
