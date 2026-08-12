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
	cookies.ClearAll(w, r)

	http.Redirect(
		w,
		r,
		"/auth/login",
		http.StatusTemporaryRedirect,
	)
}
