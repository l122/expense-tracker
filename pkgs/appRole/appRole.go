package appRole

import (
	"net/http"
)

func FromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("app_role")
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}
