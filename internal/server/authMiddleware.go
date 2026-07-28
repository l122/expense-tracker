package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/l122/expense-tracker/internal/auth"
	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/pkgs/appRole"
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

		exp, err := token.GetExpirationUnverified(accessToken)
		// _, err = token.GetExpirationUnverified(accessToken)
		if err != nil {
			redirect.ToLoginWithError(w, r, "Invalid token")
			return
		}

		ctx := token.NewContext(r.Context(), accessToken)
		if exp.Compare(time.Now()) < 1 {
			rt, err := refreshToken.FromRequest(r)
			if err != nil {
				redirect.ToLoginWithError(w, r, "No refresh token in request")
				return
			}

			tokenResp, err := auth.ExchangeRefreshToken(rt)
			if err != nil {
				// TODO: log
				redirect.ToLoginWithError(w, r, err.Error())
				return
			}

			exp := time.Unix(int64(tokenResp.ExpiresAt), 0)
			accessToken = tokenResp.AccessToken

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

			token.ToRequest(w, r, accessToken, exp)
			refreshToken.ToRequest(w, r, tokenResp.RefreshToken)
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
