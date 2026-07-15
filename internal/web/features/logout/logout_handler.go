package logout

import (
	"net/http"

	"github.com/l122/expense-tracker/pkgs/cookies"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookies.Clear(w, r)

	http.Redirect(
		w,
		r,
		"/login",
		http.StatusTemporaryRedirect,
	)
}
