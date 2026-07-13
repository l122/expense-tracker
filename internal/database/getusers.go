package database

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/l122/expense-tracker/internal/domain"
	"github.com/l122/expense-tracker/pkgs/dbhttp"
	"github.com/l122/expense-tracker/pkgs/users"
)

const (
	tableName = "profiles"
)

func (s *service) GetUsers(ctx context.Context) ([]domain.User, error) {
	var result []domain.User

	req, err := createGetUsersRequest()
	if err != nil {
		return result, err
	}

	resp, err := dbhttp.Send(ctx, req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("Error response (Status %d): %s\n", resp.StatusCode, string(bodyBytes))
		return result, err
	}

	result, err = users.FromHttpResponse(resp)
	if err != nil {
		return result, err
	}

	return result, nil
}

func createGetUsersRequest() (*http.Request, error) {
	url := os.Getenv("SUPABASE_URL") + "/rest/v1"
	endpoint := fmt.Sprintf("%s/%s?select=*", url, tableName)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return nil, err
	}

	return req, nil
}
