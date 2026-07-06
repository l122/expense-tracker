package login

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/l122/expense-tracker/internal/database"
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

	// // Fetch or Register the User in the DB
	// var dbUser domain.User
	// dbUser, err = h.db.GetUserByEmail(googleUser.Email)
	// if err != nil {
	// 	// User does not exist, let's create them
	// 	username := "usr_" + googleUser.ID
	// 	if len(username) > 15 {
	// 		username = username[:15]
	// 	}
	// 	dbUser, err = h.db.CreateUser(googleUser.Name, googleUser.Email, username, "user")
	// 	if err != nil {
	// 		http.Redirect(w, r, "/login?error=Failed+to+register+user", http.StatusTemporaryRedirect)
	// 		return
	// 	}
	// }

	shouldReturn = setSessionToken(w, r, tokenResp)
	if shouldReturn {
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func setSessionToken(w http.ResponseWriter, r *http.Request, tokenResp *TokenResponse) bool {
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		setTemporaryRedirectWithError(w, r, "SESSION_SECRET is not configured")
		return true
	}

	signedValue := signValue(tokenResp.User.ID, sessionSecret)

	// Create signed session cookie using HMAC-SHA256
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    signedValue,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
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

		setTemporaryRedirectWithError(w, r, "failed to exchange token")
		return nil, true
	}
	defer response.Body.Close()

	var tokenResp = &TokenResponse{}
	if err := json.NewDecoder(response.Body).Decode(tokenResp); err != nil {
		setTemporaryRedirectWithError(w, r, "failed to decode userinfo")
		return nil, true
	}
	return tokenResp, false
}

func getTokenRequest(r *http.Request, w http.ResponseWriter) (*http.Request, bool) {
	cookie, err := r.Cookie(verifierCookieName)
	if err != nil {
		setTemporaryRedirectWithError(w, r, "missing pkce cookie")
		return nil, true
	}
	removeAuthCookie(w)

	code := r.URL.Query().Get("code")
	if code == "" {
		setTemporaryRedirectWithError(w, r, "missing oauth code")
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

func setTemporaryRedirectWithError(w http.ResponseWriter, r *http.Request, errorMessage string) {
	errorParams := strings.Split(errorMessage, " ")
	errorMessageNormilized := strings.Join(errorParams, "+")
	http.Redirect(w, r, "/login?error="+errorMessageNormilized, http.StatusTemporaryRedirect)
}

func createTokenRequest(w http.ResponseWriter, r *http.Request, code, verifier string) (*http.Request, error) {
	payload := map[string]string{
		authCode:     code,
		codeVerifier: verifier,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		setTemporaryRedirectWithError(w, r, "json marshal failed")
		return nil, err
	}

	request, err := http.NewRequest(
		http.MethodPost,
		os.Getenv("SUPABASE_URL")+os.Getenv("SUPABASE_TOKEN_ENDPOINT"),
		bytes.NewReader(data),
	)
	if err != nil {
		setTemporaryRedirectWithError(w, r, "failed to create token exchange request")
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
