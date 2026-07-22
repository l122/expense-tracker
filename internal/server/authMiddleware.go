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
			redirect.ToLoginWithError(w, r, "no token in request")
			return
		}
		ctx := token.NewContext(r.Context(), accessToken)

		token := r.Header.Get("Authorization")
		err = validateToken(token)
		if err != nil {
			// TODO:log
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateToken(token string) error {
	// TODO: validate expiration
	// if token == "Bearer valid-token" {
	// 	return "12345"
	// }
	return nil
}
