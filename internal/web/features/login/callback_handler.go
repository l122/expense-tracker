package login

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"time"

	"github.com/l122/expense-tracker/internal/auth"
	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/pkgs/appRole"
	"github.com/l122/expense-tracker/pkgs/cookies"
	"github.com/l122/expense-tracker/pkgs/redirect"
	"github.com/l122/expense-tracker/pkgs/refreshToken"
	"github.com/l122/expense-tracker/pkgs/token"
	"github.com/l122/expense-tracker/pkgs/userid"
)

var refreshTokenDurationHours int

type CallbackHandler struct {
	view *LoginView
	db   database.Service
}

func NewCallbackHandler(view *LoginView, db database.Service) *CallbackHandler {
	duration, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_DURATION_HOURS"))
	if err != nil {
		log.Fatal("NewCallbackHandler: error parsing REFRESH_TOKEN_DURATION_HOURS:", err)
	}

	refreshTokenDurationHours = duration

	return &CallbackHandler{
		view: view,
		db:   db,
	}
}

// GET /auth/google/callback
// Supabase redirects here after Google authentication.
// Exchanges the authorization code for a session and sets a session cookie.

func (h *CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := cookies.FromRequest(r, verifierCookieName)
	if err != nil {
		redirect.ToLoginWithError(w, r, "missing pkce cookie")
		return
	}
	cookies.Clear(w, r, verifierCookieName)

	code := r.URL.Query().Get("code")
	if code == "" {
		redirect.ToLoginWithError(w, r, "missing oauth code")
		return
	}

	tokenResp, err := auth.ExchangeCodeForToken(code, cookie.Value)
	if err != nil {
		redirect.ToLoginWithError(w, r, err.Error())
		return
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

	exp := time.Unix(int64(tokenResp.ExpiresAt), 0)
	appRole.ToRequest(w, r, user.AppRole, exp)
	userid.ToRequest(w, r, user.Id, exp)
	token.ToRequest(w, r, tokenResp.AccessToken, exp)
	refreshToken.ToRequest(w, r, tokenResp.RefreshToken)

	// Patch Avatar
	_, err = h.db.UpdateAvatar(ctx, user.Id, tokenResp.User.UserMetadata.AvatarURL)
	if err != nil {
		// TODO:log
		fmt.Printf("error updating avatar url: %s", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
