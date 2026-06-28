package login

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/internal/domain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthHandler struct {
	templ *template.Template
	db    database.Service
}

func NewAuthHandler(templ *template.Template, db database.Service) *AuthHandler {
	return &AuthHandler{
		templ: templ,
		db:    db,
	}
}

// 1. Renders the Login Page
func (h *AuthHandler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	errParam := r.URL.Query().Get("error")
	data := map[string]interface{}{
		"Error": errParam,
	}
	if err := h.templ.ExecuteTemplate(w, "login.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// 2. Redirects user to Google OAuth Consent Page
func (h *AuthHandler) RedirectToGoogle(w http.ResponseWriter, r *http.Request) {
	// Generate random state to protect against CSRF
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}
	state := hex.EncodeToString(b)

	// Save state in a short-lived HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		Expires:  time.Now().Add(5 * time.Minute),
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})

	config := getOauthConfig()
	url := config.AuthCodeURL(state, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// 3. Callback URL that Google redirects to
func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Verify State
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value == "" {
		http.Redirect(w, r, "/login?error=Missing+oauth+state", http.StatusTemporaryRedirect)
		return
	}

	// Delete state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	if r.URL.Query().Get("state") != stateCookie.Value {
		http.Redirect(w, r, "/login?error=Invalid+oauth+state", http.StatusTemporaryRedirect)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Redirect(w, r, "/login?error=Missing+oauth+code", http.StatusTemporaryRedirect)
		return
	}

	config := getOauthConfig()
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		http.Redirect(w, r, "/login?error=Failed+to+exchange+token", http.StatusTemporaryRedirect)
		return
	}

	// Fetch user details from Google userinfo API
	client := config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Redirect(w, r, "login?error=Failed+to+fetch+userinfo", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		http.Redirect(w, r, "/login?error=Failed+to+decode+userinfo", http.StatusTemporaryRedirect)
		return
	}

	// Fetch or Register the User in the DB
	var dbUser domain.User
	dbUser, err = h.db.GetUserByEmail(googleUser.Email)
	if err != nil {
		// User does not exist, let's create them
		username := "usr_" + googleUser.ID
		if len(username) > 15 {
			username = username[:15]
		}
		dbUser, err = h.db.CreateUser(googleUser.Name, googleUser.Email, username, "user")
		if err != nil {
			http.Redirect(w, r, "/login?error=Failed+to+register+user", http.StatusTemporaryRedirect)
			return
		}
	}

	// Create signed session cookie using HMAC-SHA256
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = "fallback-insecure-secret" // Just in case, but warn in production
	}

	signedValue := signValue(dbUser.Email, sessionSecret)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    signedValue,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
}

// Cryptographic Cookie Signature Utilities
func signValue(value, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(value))
	signature := hex.EncodeToString(mac.Sum(nil))
	return value + "." + signature
}
