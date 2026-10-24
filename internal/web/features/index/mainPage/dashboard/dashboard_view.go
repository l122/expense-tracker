package dashboard

import (
	"html/template"
	"net/http"

	"github.com/l122/expense-tracker/internal/domain"
)

type View struct {
	templ *template.Template
}

func New(templ *template.Template) *View {
	return &View{
		templ: templ,
	}
}

func (i *View) Index(w http.ResponseWriter, categories []domain.Category) {
	if err := i.templ.ExecuteTemplate(w, "dashboard", categories); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
