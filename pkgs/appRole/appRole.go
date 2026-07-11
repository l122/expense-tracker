package appRole

import (
	"net/http"
	"time"
)

const (
	appRole = "app_role"
)

func FromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie(appRole)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func ToRequest(w http.ResponseWriter, r *http.Request, userRole string, exp time.Time) bool {
	http.SetCookie(w, &http.Cookie{
		Name:     appRole,
		Value:    userRole,
		Path:     "/",
		HttpOnly: true,
		Expires:  exp,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})

	return false
}
