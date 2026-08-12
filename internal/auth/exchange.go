package auth

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/l122/expense-tracker/pkgs/authVerifier"
	"github.com/l122/expense-tracker/pkgs/dbhttp"
	"github.com/l122/expense-tracker/pkgs/exchangeCode"
	"github.com/l122/expense-tracker/pkgs/refreshToken"
	"github.com/l122/expense-tracker/pkgs/token"
)

func ExchangeRefreshToken(rt string) (*token.TokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = refreshToken.NewContext(ctx, rt)
	response, err := dbhttp.ExchangeRefreshToken(ctx)
	if err != nil {
		// TODO: log err
		return nil, errors.New("failed to exchange token")
	}
	defer response.Body.Close()

	// Todo: FIX THIS ERROR
	var tokenResp = &token.TokenResponse{}
	if err := json.NewDecoder(response.Body).Decode(tokenResp); err != nil {
		// TODO: log err
		return nil, errors.New("failed to decode userinfo")
	}

	if !tokenResp.User.UserMetadata.EmailVerified {
		return nil, errors.New("email not verified")
	}

	return tokenResp, nil
}

func ExchangeCodeForToken(code, verifier string) (*token.TokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = exchangeCode.NewContext(ctx, code)
	ctx = authVerifier.NewContext(ctx, verifier)

	response, err := dbhttp.ExchangeCodeForToken(ctx)
	if err != nil {
		// TODO: log err
		return nil, errors.New("failed to exchange token")
	}
	defer response.Body.Close()

	var tokenResp = &token.TokenResponse{}
	if err := json.NewDecoder(response.Body).Decode(tokenResp); err != nil {
		return nil, errors.New("failed to decode userinfo")
	}

	if !tokenResp.User.UserMetadata.EmailVerified {
		return nil, errors.New("email not verified")
	}

	return tokenResp, nil
}
