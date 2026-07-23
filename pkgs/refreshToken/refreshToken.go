package refreshToken

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	refreshToken = "refresh_token"
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
