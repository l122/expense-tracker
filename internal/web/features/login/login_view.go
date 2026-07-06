package login

import (
	"html/template"
	"net/http"
)

type LoginView struct {
	templ *template.Template
}

func NewLoginView(templ *template.Template) *LoginView {
	return &LoginView{
		templ: templ,
	}
}

func (i *LoginView) Index(w http.ResponseWriter, data map[string]interface{}) {
	if err := i.templ.ExecuteTemplate(w, "login.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
