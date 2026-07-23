package login

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"time"

	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/pkgs/appRole"
	"github.com/l122/expense-tracker/pkgs/cookies"
	"github.com/l122/expense-tracker/pkgs/redirect"
	"github.com/l122/expense-tracker/pkgs/token"
	"github.com/l122/expense-tracker/pkgs/userid"
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

	user, err := h.db.GetUserByAuthId(ctx, tokenResp.User.ID)
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
	appRole.ToRequest(w, r, user.AppRole, exp)
	userid.ToRequest(w, r, user.Id, exp)
	cookies.Set(w, r, accessToken, tokenResp.AccessToken, exp)
	cookies.Set(w, r, refreshToken, tokenResp.RefreshToken, time.Now().Add(24*7*4*time.Hour))

	// Patch Avatar
	_, err = h.db.UpdateAvatar(ctx, user.Id, tokenResp.User.UserMetadata.AvatarURL)
	if err != nil {
		// TODO:log
		fmt.Printf("error updating avatar url: %s", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func exchangeCodeForToken(request *http.Request, w http.ResponseWriter, r *http.Request) (*TokenResponse, bool) {
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil || response.StatusCode > 299 || response.StatusCode < 200 {
		// TODO: log
		// body, _ := io.ReadAll(response.Body)
		// fmt.Println("STATUS:", response.Status)
		// fmt.Println("BODY:", string(body))

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
	cookies.Clear(w, r, verifierCookieName)

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
