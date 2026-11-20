package shared

import (
	"html/template"
	"net/http"
)

type NewExpenseView struct {
	templ *template.Template
}

func NewNewExpenseView(templ *template.Template) *NewExpenseView {
	return &NewExpenseView{
		templ: templ,
	}
}

func (i *NewExpenseView) Index(w http.ResponseWriter, data any) {
	if err := i.templ.ExecuteTemplate(w, "new_expense", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
