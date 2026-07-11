package admin

import (
	"html/template"
	"net/http"

	"github.com/l122/expense-tracker/internal/domain"
)

type AdminView struct {
	templ *template.Template
}

func NewAdminView(templ *template.Template) *AdminView {
	return &AdminView{
		templ: templ,
	}
}

func (i *AdminView) Index(w http.ResponseWriter, users []domain.User) {
	if err := i.templ.ExecuteTemplate(w, "admin.html", users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
