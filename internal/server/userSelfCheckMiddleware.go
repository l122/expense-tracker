package server

import (
	"fmt"
	"net/http"

	"github.com/l122/expense-tracker/pkgs/userid"
)

func userSelfCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request comes from the same user
		userIdFromCookie := userid.FromRequestOrPanic(r)
		userIdFromUrl, err := userid.FromUrlPath(r)
		if err != nil {
			// TODO: log
			fmt.Printf("Error: %v\n", err)
			http.Error(w, "Invalid or missing user id in URL path", http.StatusBadRequest)

			return
		}
		if userIdFromCookie != userIdFromUrl {
			// TODO: log and redirect to an error page
			// return 403
			http.Error(w, "User Ids in cookie and in URL do not match", http.StatusForbidden)

			return
		}

		next.ServeHTTP(w, r)
	})
}
