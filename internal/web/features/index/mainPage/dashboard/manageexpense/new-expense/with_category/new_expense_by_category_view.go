package with_category

import (
	"html/template"
	"net/http"
)

type NewExpenseByCategoryView struct {
	templ *template.Template
}

func NewNewExpenseByCategoryView(templ *template.Template) *NewExpenseByCategoryView {
	return &NewExpenseByCategoryView{
		templ: templ,
	}
}

func (i *NewExpenseByCategoryView) Index(w http.ResponseWriter, data any) {
	if err := i.templ.ExecuteTemplate(w, "new_expense_by_category", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
