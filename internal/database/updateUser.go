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
	"github.com/l122/expense-tracker/pkgs/token"
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

	url := os.Getenv("SUPABASE_URL") + "/rest/v1"
	endpoint := fmt.Sprintf("%s/%s?id=eq.%s", url, tableName, userId)
	payload := map[string]interface{}{
		"enabled": enable,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return emptyUser, err
	}

	req, err := http.NewRequest("PATCH", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Failed to patch user: %v\n", err)
		return emptyUser, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")
	req.Header.Set("apikey", os.Getenv("SUPABASE_PUBLISHABLE_KEY"))
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("HTTP request failed: %v\n", err)
		return emptyUser, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("Error response (Status %d): %s\n", resp.StatusCode, string(bodyBytes))
		return emptyUser, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	fmt.Println("BODY:", string(bodyBytes))
	if err != nil {
		fmt.Printf("Failed to read body: %v\n", err)
		return emptyUser, err
	}

	var users []domain.User
	if err := json.Unmarshal(bodyBytes, &users); err != nil || len(users) == 0 {
		fmt.Printf("Failed to parse JSON: %v\n", err)
		return emptyUser, err
	}

	return users[0], nil
}
