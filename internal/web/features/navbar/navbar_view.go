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

func (i *View) Index(w http.ResponseWriter, role string) {
	if err := i.templ.ExecuteTemplate(w, "navbar", role); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
