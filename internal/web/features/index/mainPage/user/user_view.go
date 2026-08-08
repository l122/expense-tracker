package user

import (
	"html/template"
	"net/http"

	"github.com/l122/expense-tracker/internal/domain"
)

type UserView struct {
	templ *template.Template
}

func NewView(templ *template.Template) *UserView {
	return &UserView{
		templ: templ,
	}
}

func (i *UserView) Index(w http.ResponseWriter, user domain.User) {
	if err := i.templ.ExecuteTemplate(w, "user", user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
