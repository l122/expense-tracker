package server

import (
	"net/http"

	"github.com/l122/expense-tracker/pkgs/redirect"
	"github.com/l122/expense-tracker/pkgs/userid"
)

func userIdCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := userid.FromRequest(r)
		if err != nil {
			redirect.ToLoginWithError(w, r, "no user id in request")
			return
		}

		next.ServeHTTP(w, r)
	})
}
