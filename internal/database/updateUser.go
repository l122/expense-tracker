package database

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/l122/expense-tracker/internal/domain"
	"github.com/l122/expense-tracker/pkgs/myhttp"
	"github.com/l122/expense-tracker/pkgs/token"
	"github.com/l122/expense-tracker/pkgs/users"
)

func (s *service) EnableUser(ctx context.Context, userId string) (domain.User, error) {
	return patchEnable(ctx, userId, true)
}

func (s *service) DisableUser(ctx context.Context, userId string) (domain.User, error) {
	return patchEnable(ctx, userId, false)
}

func patchEnable(ctx context.Context, userId string, enable bool) (domain.User, error) {
	var emptyUser domain.User
	token, ok := token.FromContext(ctx)
	if !ok {
		return emptyUser, errors.New("Failed to extract token")
	}

	req, err := createRequest(userId, enable, token)
	if err != nil {
		return emptyUser, err
	}

	resp, err := myhttp.Send(req)
	if err != nil {
		return emptyUser, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("Error response (Status %d): %s\n", resp.StatusCode, string(bodyBytes))
		return emptyUser, err
	}

	users, err := users.FromHttpResponse(resp, emptyUser)
	if err != nil {
		return emptyUser, err
	}

	if len(users) == 0 {
		fmt.Println("error updating user")
		return emptyUser, errors.New("error updating user")
	}

	return users[0], nil
}

func createRequest(userId string, enable bool, token string) (*http.Request, error) {
	url := os.Getenv("SUPABASE_URL") + "/rest/v1"
	endpoint := fmt.Sprintf("%s/%s?id=eq.%s", url, tableName, userId)
	payload := map[string]interface{}{
		"enabled": enable,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return nil, err
	}

	req, err := http.NewRequest("PATCH", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Failed to patch user: %v\n", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")
	req.Header.Set("apikey", os.Getenv("SUPABASE_PUBLISHABLE_KEY"))
	req.Header.Set("Authorization", "Bearer "+token)
	return req, nil
}
