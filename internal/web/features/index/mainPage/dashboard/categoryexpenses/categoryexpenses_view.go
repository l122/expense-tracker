package categoryexpenses

import (
	"html/template"
	"net/http"

	"github.com/l122/expense-tracker/internal/web/features/index/mainPage/dashboard/categoryexpenses/models"
)

type View struct {
	templ *template.Template
}

func New(templ *template.Template) *View {
	return &View{
		templ: templ,
	}
}

func (i *View) Index(w http.ResponseWriter, data models.CategoryExpensesData) {
	if err := i.templ.ExecuteTemplate(w, "categoryexpenses", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
