package database

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/l122/expense-tracker/internal/domain"
	"github.com/l122/expense-tracker/pkgs/token"
)

const (
	tableName = "profiles"
)

func (s *service) GetUsers(ctx context.Context) ([]domain.User, error) {
	token, ok := token.FromContext(ctx)
	if ok {
		fmt.Println(token)
		// request.Set("Bearer", token)
	}

	url := os.Getenv("SUPABASE_URL") + "/rest/v1"
	endpoint := fmt.Sprintf("%s/%s?select=*", url, tableName)

	var result []domain.User
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return result, err
	}

	req.Header.Set("apikey", os.Getenv("SUPABASE_PUBLISHABLE_KEY"))
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("HTTP request failed: %v\n", err)
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("Error response (Status %d): %s\n", resp.StatusCode, string(bodyBytes))
		return result, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read body: %v\n", err)
		return result, err
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		fmt.Printf("Failed to parse JSON: %v\n", err)
		return result, err
	}

	return result, nil
}
