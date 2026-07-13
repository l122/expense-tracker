package userid

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	idKey = "user_id"
)

func FromRequest(r *http.Request) (int, error) {
	cookie, err := r.Cookie(idKey)
	if err != nil {
		return 0, err
	}

	res, err := strconv.Atoi(cookie.Value)
	if err != nil {
		fmt.Println("Error parsing string:", err)
		return 0, err
	}

	return res, nil
}

func ToRequest(w http.ResponseWriter, r *http.Request, userId int, exp time.Time) bool {
	http.SetCookie(w, &http.Cookie{
		Name:     idKey,
		Value:    strconv.Itoa(userId),
		Path:     "/",
		HttpOnly: true,
		Expires:  exp,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})

	return false
}
