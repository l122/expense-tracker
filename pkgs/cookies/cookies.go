package cookies

import (
	"net/http"
	"time"
)

func Clear(w http.ResponseWriter, r *http.Request) {
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
