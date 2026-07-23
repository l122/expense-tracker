package token

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type key int

const (
	accessToken              = "access_token"
	missingExpClaimError     = "exp claim missing from payload"
	missingSubClaimError     = "sub claim missing from payload"
	tokenKey             key = 0
)

func FromRequestHeader(r *http.Request) string {
	return r.Header.Get("Authorization")
}

func FromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie(accessToken)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func FromRequestOrPanic(r *http.Request) string {
	accessToken, err := FromRequest(r)
	if err != nil {
		panic("no token in request")
	}

	return accessToken
}

func NewContext(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}

func FromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(tokenKey).(string)
	return token, ok
}

func GetExpirationUnverified(tokenString string) (time.Time, error) {
	parser := jwt.NewParser()
	claims := jwt.MapClaims{}

	// Parse unverified skips signature verification
	_, _, err := parser.ParseUnverified(tokenString, &claims)
	if err != nil {
		return time.Time{}, err
	}

	// Extract the expiration time
	if exp, err := claims.GetExpirationTime(); err == nil && exp != nil {
		return exp.Time, nil
	}

	return time.Time{}, fmt.Errorf(missingExpClaimError)
}

func GetUserId(tokenString string) (string, error) {
	parser := jwt.NewParser()
	claims := jwt.MapClaims{}

	// Parse unverified skips signature verification
	_, _, err := parser.ParseUnverified(tokenString, &claims)
	if err != nil {
		return "", err
	}

	if sub, err := claims.GetSubject(); err == nil && sub != "" {
		return sub, nil
	}

	return "", fmt.Errorf(missingSubClaimError)
}
