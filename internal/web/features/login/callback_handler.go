package login

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"time"

	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/pkgs/appRole"
	"github.com/l122/expense-tracker/pkgs/redirect"
	"github.com/l122/expense-tracker/pkgs/token"
)

type CallbackHandler struct {
	view *LoginView
	db   database.Service
}

func NewCallbackHandler(view *LoginView, db database.Service) *CallbackHandler {
	return &CallbackHandler{
		view: view,
		db:   db,
	}
}

// GET /auth/google/callback
// Supabase redirects here after Google authentication.
// Exchanges the authorization code for a session and sets a session cookie.

func (h *CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request, shouldReturn := getTokenRequest(r, w)
	if shouldReturn {
		return
	}

	tokenResp, shouldReturn := exchangeCodeForToken(request, w, r)
	if shouldReturn {
		return
	}

	if !tokenResp.User.UserMetadata.EmailVerified {
		redirect.ToLoginWithError(w, r, "email not verified")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = token.NewContext(ctx, tokenResp.AccessToken)

	user, err := h.db.GetUserById(ctx, tokenResp.User.ID)
	if err != nil {
		redirect.ToLoginWithError(w, r, "user not found")
		return
	}

	if !user.Enabled {
		redirect.ToLoginWithError(w, r, "user not enabled")
		return
	}

	exp, err := token.GetExpirationUnverified(tokenResp.AccessToken)
	if err != nil {
		redirect.ToLoginWithError(w, r, err.Error())
		return
	}
	shouldReturn = appRole.ToRequest(w, r, user.AppRole, exp)
	if shouldReturn {
		return
	}

	shouldReturn = setSessionTokens(w, r, tokenResp)
	if shouldReturn {
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func setSessionTokens(w http.ResponseWriter, r *http.Request, tokenResp *TokenResponse) bool {
	exp, err := token.GetExpirationUnverified(tokenResp.AccessToken)
	if err != nil {
		redirect.ToLoginWithError(w, r, err.Error())
	}

	http.SetCookie(w, &http.Cookie{
		Name:     accessToken,
		Value:    tokenResp.AccessToken,
		Path:     "/",
		Expires:  exp,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     refreshToken,
		Value:    tokenResp.RefreshToken,
		Path:     "/",
		Expires:  time.Now().Add(24 * 7 * 4 * time.Hour), // todo: move to configs
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
	return false
}

func exchangeCodeForToken(request *http.Request, w http.ResponseWriter, r *http.Request) (*TokenResponse, bool) {
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil || response.StatusCode > 299 || response.StatusCode < 200 {
		// TODO: log
		body, _ := io.ReadAll(response.Body)
		fmt.Println("STATUS:", response.Status)
		fmt.Println("BODY:", string(body))

		redirect.ToLoginWithError(w, r, "failed to exchange token")
		return nil, true
	}
	defer response.Body.Close()

	var tokenResp = &TokenResponse{}
	if err := json.NewDecoder(response.Body).Decode(tokenResp); err != nil {
		redirect.ToLoginWithError(w, r, "failed to decode userinfo")
		return nil, true
	}
	return tokenResp, false
}

func getTokenRequest(r *http.Request, w http.ResponseWriter) (*http.Request, bool) {
	cookie, err := r.Cookie(verifierCookieName)
	if err != nil {
		redirect.ToLoginWithError(w, r, "missing pkce cookie")
		return nil, true
	}
	removeAuthCookie(w)

	code := r.URL.Query().Get("code")
	if code == "" {
		redirect.ToLoginWithError(w, r, "missing oauth code")
		return nil, true
	}

	request, err := createTokenRequest(w, r, code, cookie.Value)
	if err != nil {
		return nil, true
	}
	return request, false
}

// Cryptographic Cookie Signature Utilities
func signValue(value, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(value))
	signature := hex.EncodeToString(mac.Sum(nil))
	return value + "." + signature
}

func createTokenRequest(w http.ResponseWriter, r *http.Request, code, verifier string) (*http.Request, error) {
	payload := map[string]string{
		authCode:     code,
		codeVerifier: verifier,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		redirect.ToLoginWithError(w, r, "json marshal failed")
		return nil, err
	}

	request, err := http.NewRequest(
		http.MethodPost,
		os.Getenv("SUPABASE_URL")+os.Getenv("SUPABASE_TOKEN_ENDPOINT"),
		bytes.NewReader(data),
	)
	if err != nil {
		redirect.ToLoginWithError(w, r, "failed to create token exchange request")
		return nil, err
	}

	request.Header.Set(contentType, applicationJson)
	request.Header.Set(apiKey, os.Getenv("SUPABASE_PUBLISHABLE_KEY"))

	return request, nil
}

func removeAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     verifierCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
}
