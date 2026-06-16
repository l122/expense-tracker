package about

import (
	"html/template"
	"net/http"

	"github.com/l122/expense-tracker/internal/config"
)

type AboutView struct {
	templ  *template.Template
	config *config.Config
}

func NewAboutView(templ *template.Template, conf *config.Config) *AboutView {
	return &AboutView{
		templ:  templ,
		config: conf,
	}
}

func (i *AboutView) Index(w http.ResponseWriter) {
	if err := i.templ.ExecuteTemplate(w, "about.html", i.config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
