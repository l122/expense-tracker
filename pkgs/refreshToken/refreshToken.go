package refreshToken

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type key int

const (
	refreshToken     = "refresh_token"
	tokenKey     key = 0
)

var refreshTokenDurationHours int

func init() {
	duration, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_DURATION_HOURS"))
	if err != nil {
		log.Fatal("NewCallbackHandler: error parsing REFRESH_TOKEN_DURATION_HOURS:", err)
	}

	refreshTokenDurationHours = duration
}

func ToRequest(w http.ResponseWriter, r *http.Request, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshToken,
		Value:    value,
		Path:     "/",
		Expires:  time.Now().Add(time.Duration(refreshTokenDurationHours) * time.Hour),
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
}

func FromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie(refreshToken)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func NewContext(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}

func FromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(tokenKey).(string)
	return token, ok
}
