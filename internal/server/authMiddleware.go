package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/l122/expense-tracker/pkgs/dbhttp"
	"github.com/l122/expense-tracker/pkgs/redirect"
	"github.com/l122/expense-tracker/pkgs/refreshToken"
	"github.com/l122/expense-tracker/pkgs/token"
)

func authMiddleware(next http.Handler) http.Handler {
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

			// TODO: remove after debug
			body, _ := io.ReadAll(response.Body)
			fmt.Println("STATUS:", response.Status)
			fmt.Println("BODY:", string(body))

			// Todo: FIX THIS ERROR
			var tokenResp = &RefreshTokenResponse{}
			if err := json.NewDecoder(response.Body).Decode(tokenResp); err != nil {
				redirect.ToLoginWithError(w, r, "failed to decode userinfo")
				return
			}

			exp, err := token.GetExpirationUnverified(accessToken)
			if err != nil {
				redirect.ToLoginWithError(w, r, "Invalid token")
				return
			}
			token.ToRequest(w, r, tokenResp.AccessToken, exp)
			refreshToken.ToRequest(w, r, tokenResp.RefreshToken)

			// Patch Avatar
			// _, err = h.db.UpdateAvatar(ctx, user.Id, tokenResp.User.UserMetadata.AvatarURL)
			// if err != nil {
			// 	// TODO:log
			// 	fmt.Printf("error updating avatar url: %s", err)
			// }

			refreshToken.ToRequest(w, r, rt)

			// exp, err := token.GetExpirationUnverified(newAccessToken)
			// if err != nil {
			// 	redirect.ToLoginWithError(w, r, err.Error())
			// 	return
			// }
			// token.ToRequest(w, r, newAccessToken, exp)
		}

		ctx := token.NewContext(r.Context(), accessToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExpiresAt    int    `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		ID               string    `json:"id"`
		Aud              string    `json:"aud"`
		Role             string    `json:"role"`
		Email            string    `json:"email"`
		EmailConfirmedAt time.Time `json:"email_confirmed_at"`
		Phone            string    `json:"phone"`
		ConfirmedAt      time.Time `json:"confirmed_at"`
		LastSignInAt     time.Time `json:"last_sign_in_at"`
		AppMetadata      struct {
			Provider  string   `json:"provider"`
			Providers []string `json:"providers"`
		} `json:"app_metadata"`
		UserMetadata struct {
			AvatarURL     string `json:"avatar_url"`
			Email         string `json:"email"`
			EmailVerified bool   `json:"email_verified"`
			FullName      string `json:"full_name"`
			Iss           string `json:"iss"`
			Name          string `json:"name"`
			PhoneVerified bool   `json:"phone_verified"`
			Picture       string `json:"picture"`
			ProviderID    string `json:"provider_id"`
			Sub           string `json:"sub"`
		} `json:"user_metadata"`
		Identities []struct {
			IdentityID   string `json:"identity_id"`
			ID           string `json:"id"`
			UserID       string `json:"user_id"`
			IdentityData struct {
				AvatarURL     string `json:"avatar_url"`
				Email         string `json:"email"`
				EmailVerified bool   `json:"email_verified"`
				FullName      string `json:"full_name"`
				Iss           string `json:"iss"`
				Name          string `json:"name"`
				PhoneVerified bool   `json:"phone_verified"`
				Picture       string `json:"picture"`
				ProviderID    string `json:"provider_id"`
				Sub           string `json:"sub"`
			} `json:"identity_data"`
			Provider     string    `json:"provider"`
			LastSignInAt time.Time `json:"last_sign_in_at"`
			CreatedAt    time.Time `json:"created_at"`
			UpdatedAt    time.Time `json:"updated_at"`
			Email        string    `json:"email"`
		} `json:"identities"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		IsAnonymous bool      `json:"is_anonymous"`
	} `json:"user"`
}
