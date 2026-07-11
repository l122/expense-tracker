package navbar

import (
	"net/http"
)

type NavbarHandler struct {
	http.Handler

	view *View
}

func NewHandler(view *View) *NavbarHandler {
	return &NavbarHandler{
		view: view,
	}
}

func (t *NavbarHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	t.view.Index(w)
}
