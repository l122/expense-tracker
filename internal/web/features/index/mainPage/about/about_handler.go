package about

import (
	"net/http"
)

type AboutHandler struct {
	http.Handler

	aboutView *AboutView
}

func NewAboutHandler(aboutView *AboutView) *AboutHandler {
	return &AboutHandler{
		aboutView: aboutView,
	}
}

func (t *AboutHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	t.aboutView.Index(w)
}
