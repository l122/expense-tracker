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

func FromRequestOrPanic(r *http.Request) int {
	userId, err := FromRequest(r)
	if err != nil {
		panic("no user id in request")
	}

	return userId
}

func FromUrlPath(r *http.Request) (int, error) {
	userId := r.PathValue("id")

	res, err := strconv.Atoi(userId)
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
