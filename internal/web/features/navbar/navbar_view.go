package navbar

import (
	"html/template"
	"net/http"

	"github.com/l122/expense-tracker/internal/domain"
)

type View struct {
	templ *template.Template
}

func NewView(templ *template.Template) *View {
	return &View{
		templ: templ,
	}
}

func (i *View) Index(w http.ResponseWriter, user domain.User) {
	if err := i.templ.ExecuteTemplate(w, "navbar", user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
