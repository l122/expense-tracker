package dbhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/l122/expense-tracker/pkgs/authVerifier"
	"github.com/l122/expense-tracker/pkgs/exchangeCode"
	"github.com/l122/expense-tracker/pkgs/refreshToken"
	"github.com/l122/expense-tracker/pkgs/token"
)

const (
	authCode        = "auth_code"
	codeVerifier    = "code_verifier"
	contentType     = "Content-Type"
	applicationJson = "application/json"
	apiKey          = "apikey"
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

func SendWithRoleKey(ctx context.Context, req *http.Request) (*http.Response, error) {
	req.Header.Set("apikey", os.Getenv("SUPABASE_API_KEY"))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_API_KEY"))

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

func ExchangeCodeForToken(ctx context.Context) (*http.Response, error) {
	code, ok := exchangeCode.FromContext(ctx)
	if !ok {
		return nil, errors.New("Failed to extract code from context")
	}

	verifier, ok := authVerifier.FromContext(ctx)
	if !ok {
		return nil, errors.New("Failed to extract verifier from context")
	}

	client := &http.Client{}
	request, err := createTokenRequest(code, verifier)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil || response.StatusCode > 299 || response.StatusCode < 200 {
		// TODO: log
		// body, _ := io.ReadAll(response.Body)
		// fmt.Println("STATUS:", response.Status)
		// fmt.Println("BODY:", string(body))
		return nil, err
	}

	return response, nil
}

func createTokenRequest(code, verifier string) (*http.Request, error) {
	payload := map[string]string{
		authCode:     code,
		codeVerifier: verifier,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(
		http.MethodPost,
		os.Getenv("SUPABASE_URL")+os.Getenv("SUPABASE_TOKEN_ENDPOINT"),
		bytes.NewReader(data),
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set(contentType, applicationJson)
	request.Header.Set(apiKey, os.Getenv("SUPABASE_PUBLISHABLE_KEY"))

	return request, nil
}
