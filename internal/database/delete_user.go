package database

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/l122/expense-tracker/pkgs/dbhttp"
)

func (s *service) DeleteUser(ctx context.Context, userId int) error {
	user, err := s.GetUserById(ctx, userId)
	if err != nil {
		// TODO: log
		fmt.Printf("User not found: %v\n", err)
		return err
	}

	req, err := createDeleteRequest(user.AuthId)
	if err != nil {
		// TODO: log
		fmt.Printf("Failed to create delete request: %v\n", err)
		return err
	}

	resp, err := dbhttp.SendWithRoleKey(ctx, req)
	if err != nil {
		// TODO: log
		fmt.Printf("Failed to delete user: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("Error response (Status %d): %s\n", resp.StatusCode, string(bodyBytes))
		return err
	}

	return nil
}

func createDeleteRequest(auth_id string) (*http.Request, error) {

	url := os.Getenv("SUPABASE_URL")
	endpoint := fmt.Sprintf("%s/auth/v1/admin/users/%s", url, auth_id)

	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		fmt.Printf("Failed to delete user: %v\n", err)
		return nil, err
	}

	return req, nil
}
