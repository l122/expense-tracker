package database

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/l122/expense-tracker/internal/domain"
	"github.com/l122/expense-tracker/pkgs/categories"
	"github.com/l122/expense-tracker/pkgs/dbhttp"
	"github.com/l122/expense-tracker/pkgs/httperrors"
)

func (s *service) GetCategoryByIdForUser(ctx context.Context, categoryId int, userId int) (domain.Category, error) {
	emptyCategory := domain.Category{}
	req, err := createGetCategoryByIdForUserRequest(categoryId, userId)
	if err != nil {
		return emptyCategory, err
	}

	resp, err := dbhttp.Send(ctx, req)
	if err != nil {
		return emptyCategory, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("Error response (Status %d): %s\n", resp.StatusCode, string(bodyBytes))
		return emptyCategory, fmt.Errorf("Error response (Status %d): %s\n", resp.StatusCode, string(bodyBytes))
	}

	fetchedData, err := categories.FromHttpResponse(resp)
	if err != nil {
		return emptyCategory, err
	}

	if len(fetchedData) == 0 {
		return emptyCategory, httperrors.ErrNotFound
	}

	return fetchedData[0], nil
}

func createGetCategoryByIdForUserRequest(categoryId int, userId int) (*http.Request, error) {
	url := os.Getenv("SUPABASE_URL") + "/rest/v1"
	endpoint := fmt.Sprintf("%s/%s?select=*&id=eq.%d&profile_id=eq.%d", url, categoriesTable, categoryId, userId)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return nil, err
	}
	return req, nil
}
