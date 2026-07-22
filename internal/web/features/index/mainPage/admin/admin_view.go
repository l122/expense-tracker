package admin

import (
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
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

func (i *AdminView) Index(w http.ResponseWriter, r *http.Request, users []domain.User) {
	data := map[string]any{
		csrf.TemplateTag: csrf.TemplateField(r),
		"Items":          users,
	}

	if err := i.templ.ExecuteTemplate(w, "admin.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
