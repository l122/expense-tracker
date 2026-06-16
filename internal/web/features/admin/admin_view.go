package admin

import (
	"html/template"
	"net/http"
)

type AdminView struct {
	templ *template.Template
}

type User struct {
	Index    int
	Name     string
	Email    string
	Username string
	Role     string
}

func NewAdminView(templ *template.Template) *AdminView {
	return &AdminView{
		templ: templ,
	}
}

func (i *AdminView) Index(w http.ResponseWriter, users []User) {
	if err := i.templ.ExecuteTemplate(w, "admin.html", users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
