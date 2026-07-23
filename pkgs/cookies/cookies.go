package cookies

import (
	"net/http"
	"time"
)

func ClearAll(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()

	for _, cookie := range cookies {
		expiredCookie := &http.Cookie{
			Name:     cookie.Name,
			Value:    "",
			Path:     "/",             // Must match the original path (usually "/")
			MaxAge:   -1,              // -1 forces the browser to delete the cookie immediately
			Expires:  time.Unix(0, 0), // Legacy fallback support for older browsers
			HttpOnly: true,
		}

		http.SetCookie(w, expiredCookie)
	}
}

func Clear(w http.ResponseWriter, r *http.Request, name string) {
	expiredCookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",             // Must match the original path (usually "/")
		MaxAge:   -1,              // -1 forces the browser to delete the cookie immediately
		Expires:  time.Unix(0, 0), // Legacy fallback support for older browsers
		HttpOnly: true,
	}
	http.SetCookie(w, expiredCookie)
}

func Set(w http.ResponseWriter, r *http.Request, name, value string, exp time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  exp,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
}
