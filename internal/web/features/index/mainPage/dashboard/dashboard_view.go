package dashboard

import (
	"html/template"
	"net/http"
)

type View struct {
	templ *template.Template
}

func New(templ *template.Template) *View {
	return &View{
		templ: templ,
	}
}

func (i *View) Index(w http.ResponseWriter) {
	if err := i.templ.ExecuteTemplate(w, "dashboard", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
