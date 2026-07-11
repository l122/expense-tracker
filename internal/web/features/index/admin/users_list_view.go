package admin

import (
	"html/template"
	"net/http"

	"github.com/l122/expense-tracker/internal/domain"
)

type UsersListView struct {
	templ *template.Template
}

func NewUsersListView(templ *template.Template) *UsersListView {
	return &UsersListView{
		templ: templ,
	}
}

func (i *UsersListView) View(w http.ResponseWriter, users []domain.User) {
	if err := i.templ.ExecuteTemplate(w, "admin_container", users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
