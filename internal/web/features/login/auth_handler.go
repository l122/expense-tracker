package login

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/url"
	"os"

	"github.com/l122/expense-tracker/internal/database"
)

type AuthHandler struct {
	view *LoginView
	db   database.Service
}

func NewAuthHandler(view *LoginView, db database.Service) *AuthHandler {
	return &AuthHandler{
		view: view,
		db:   db,
	}
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	verifier, _ := generateCodeVerifier()
	http.SetCookie(w, createAuthCookie(r, verifier))
	http.Redirect(
		w,
		r,
		createAuthURL(createAuthParameters(generateCodeChallenge(verifier))),
		http.StatusTemporaryRedirect,
	)
}

func generateCodeVerifier() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func generateCodeChallenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func createAuthCookie(r *http.Request, verifier string) *http.Cookie {
	return &http.Cookie{
		Name:     verifierCookieName,
		Value:    verifier,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	}
}

func createAuthParameters(challenge string) url.Values {
	params := url.Values{}
	params.Set("provider", "google")
	params.Set("redirect_to", os.Getenv("GOOGLE_REDIRECT_URL"))
	params.Set("code_challenge", challenge)
	params.Set("code_challenge_method", "S256")
	return params
}

func createAuthURL(params url.Values) string {
	return os.Getenv("SUPABASE_URL") +
		os.Getenv("SUPABASE_AUTHORIZE_ENDPOINT") + // "/auth/v1/authorize?"
		params.Encode()
}
