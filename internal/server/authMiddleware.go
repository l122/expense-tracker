package server

import (
	"net/http"

	"github.com/l122/expense-tracker/pkgs/redirect"
	"github.com/l122/expense-tracker/pkgs/token"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := token.FromRequest(r)
		if err != nil {
			// TODO:log
			redirect.ToLoginWithError(w, r, "No access token in request")
			return
		}

		// exp, err := token.GetExpirationUnverified(accessToken)
		// if err != nil {
		// 	redirect.ToLoginWithError(w, r, "Invalid token")
		// 	return
		// }

		// if exp.Compare(time.Now()) < 1 {
		// 	// TODO: exchange token

		// 	if err != nil {
		// 		redirect.ToLoginWithError(w, r, "")
		// 		return
		// 	}
		// }

		ctx := token.NewContext(r.Context(), accessToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
