package index

import (
	"html/template"
	"net/http"

	"github.com/l122/expense-tracker/internal/config"
)

type IndexView struct {
	templ  *template.Template
	config *config.Config
}

func NewIndexView(templ *template.Template, conf *config.Config) *IndexView {
	return &IndexView{
		templ:  templ,
		config: conf,
	}
}

func (i *IndexView) Index(w http.ResponseWriter) {
	if err := i.templ.ExecuteTemplate(w, "index.html", i.config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
