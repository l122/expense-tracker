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

var users = []User{
	{
		Index:    1,
		Name:     "Alice Johnson",
		Email:    "alice.johnson@example.com",
		Username: "usr_982341",
		Role:     "admin",
	},
	{
		Index:    2,
		Name:     "Bob Smith",
		Email:    "bob.smith@provider.net",
		Username: "usr_441209",
		Role:     "user",
	},
	{
		Index:    3,
		Name:     "Charlie Davis",
		Email:    "charlie.d@service.org",
		Username: "usr_112233",
		Role:     "user",
	},
}

func NewAdminView(templ *template.Template) *AdminView {
	return &AdminView{
		templ: templ,
	}
}

func (i *AdminView) Index(w http.ResponseWriter) {
	if err := i.templ.ExecuteTemplate(w, "admin.html", users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
