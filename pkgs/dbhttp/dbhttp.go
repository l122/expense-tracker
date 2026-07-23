package dbhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/l122/expense-tracker/pkgs/refreshToken"
	"github.com/l122/expense-tracker/pkgs/token"
)

func Send(ctx context.Context, req *http.Request) (*http.Response, error) {
	token, ok := token.FromContext(ctx)
	if !ok {
		return nil, errors.New("Failed to extract token")
	}

	req.Header.Set("apikey", os.Getenv("SUPABASE_PUBLISHABLE_KEY"))
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("HTTP request failed: %v\n", err)
		return nil, err
	}
	return resp, nil
}

func ExchangeRefreshToken(ctx context.Context) (*http.Response, error) {
	rt, ok := refreshToken.FromContext(ctx)
	if !ok {
		return nil, errors.New("Failed to extract refresh token")
	}

	payload := map[string]string{
		"refresh_token": rt,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// https://<your-project-id>.supabase.co/auth/v1/token?grant_type=refresh_token
	request, err := http.NewRequest(
		http.MethodPost,
		os.Getenv("SUPABASE_URL")+os.Getenv("SUPABASE_REFRESH_TOKEN_ENDPOINT"),
		bytes.NewReader(data),
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("apikey", os.Getenv("SUPABASE_PUBLISHABLE_KEY"))

	client := &http.Client{}
	response, err := client.Do(request)

	return response, nil
}
