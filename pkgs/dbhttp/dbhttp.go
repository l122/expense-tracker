package dbhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

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
