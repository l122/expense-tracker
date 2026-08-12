package server

import (
	"net/http"

	"github.com/l122/expense-tracker/pkgs/appRole"
	"github.com/l122/expense-tracker/pkgs/redirect"
)

func adminRoleCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check role
		appRole, err := appRole.FromRequest(r)
		if err != nil {
			redirect.ToLoginWithError(w, r, "no AppRole in request")
			return
		}

		if appRole != "admin" {
			// TODO: log and redirect to an error page
			return
		}
		next.ServeHTTP(w, r)
	})
}
