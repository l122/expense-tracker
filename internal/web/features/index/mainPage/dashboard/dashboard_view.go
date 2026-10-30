package dashboard

import (
	"html/template"
	"net/http"

	"github.com/l122/expense-tracker/internal/web/features/index/mainPage/dashboard/models"
)

type View struct {
	templ *template.Template
}

func New(templ *template.Template) *View {
	return &View{
		templ: templ,
	}
}

func (i *View) Index(w http.ResponseWriter, categoriesList models.CategoriesList) {
	if err := i.templ.ExecuteTemplate(w, "dashboard", categoriesList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
