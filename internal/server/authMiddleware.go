package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/pkgs/appRole"
	"github.com/l122/expense-tracker/pkgs/dbhttp"
	"github.com/l122/expense-tracker/pkgs/redirect"
	"github.com/l122/expense-tracker/pkgs/refreshToken"
	"github.com/l122/expense-tracker/pkgs/token"
	"github.com/l122/expense-tracker/pkgs/userid"
)

func authMiddleware(next http.Handler, db database.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := token.FromRequest(r)
		if err != nil {
			// TODO:log
			redirect.ToLoginWithError(w, r, "No access token in request")
			return
		}

		// exp, err := token.GetExpirationUnverified(accessToken)
		_, err = token.GetExpirationUnverified(accessToken)
		if err != nil {
			redirect.ToLoginWithError(w, r, "Invalid token")
			return
		}

		ctx := token.NewContext(r.Context(), accessToken)
		if true {
			// if exp.Compare(time.Now()) < 1 {
			rt, err := refreshToken.FromRequest(r)
			if err != nil {
				redirect.ToLoginWithError(w, r, "No refresh token in request")
				return
			}

			// Exchange token
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			ctx = refreshToken.NewContext(ctx, rt)
			response, err := dbhttp.ExchangeRefreshToken(ctx)
			if err != nil || response.StatusCode > 299 || response.StatusCode < 200 {
				// TODO: log
				body, _ := io.ReadAll(response.Body)
				fmt.Println("STATUS:", response.Status)
				fmt.Println("BODY:", string(body))

				redirect.ToLoginWithError(w, r, "failed to exchange token")
				return
			}
			defer response.Body.Close()

			// Todo: FIX THIS ERROR
			var tokenResp = &RefreshTokenResponse{}
			if err := json.NewDecoder(response.Body).Decode(tokenResp); err != nil {
				// TODO: log
				body, _ := io.ReadAll(response.Body)
				fmt.Println("STATUS:", response.Status)
				fmt.Println("BODY:", string(body))

				redirect.ToLoginWithError(w, r, "failed to decode userinfo")
				return
			}

			exp := time.Unix(int64(tokenResp.ExpiresAt), 0)
			accessToken = tokenResp.AccessToken
			token.ToRequest(w, r, accessToken, exp)
			refreshToken.ToRequest(w, r, tokenResp.RefreshToken)

			ctx = token.NewContext(r.Context(), accessToken)
			user, err := db.GetUserByAuthId(ctx, tokenResp.User.ID)
			if err != nil {
				redirect.ToLoginWithError(w, r, "user not found")
				return
			}

			if !user.Enabled {
				redirect.ToLoginWithError(w, r, "user not enabled")
				return
			}

			appRole.ToRequest(w, r, user.AppRole, exp)
			userid.ToRequest(w, r, user.Id, exp)

			// Patch Avatar
			_, err = db.UpdateAvatar(ctx, user.Id, tokenResp.User.UserMetadata.AvatarURL)
			if err != nil {
				// TODO:log
				fmt.Printf("error updating avatar url: %s", err)
			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresAt    int    `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		ID           string `json:"id"`
		Email        string `json:"email"`
		UserMetadata struct {
			AvatarURL string `json:"avatar_url"`
			FullName  string `json:"full_name"`
		} `json:"user_metadata"`
	} `json:"user"`
}
