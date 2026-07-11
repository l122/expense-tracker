package navbar

import (
	"net/http"

	"github.com/l122/expense-tracker/pkgs/appRole"
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

func (t *NavbarHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	appRole, err := appRole.FromRequest(r)
	if err != nil {
		return
	}

	t.view.Index(w, appRole)
}
