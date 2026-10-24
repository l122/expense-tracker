package database

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/l122/expense-tracker/internal/domain"
	"github.com/l122/expense-tracker/pkgs/dbhttp"
)

const (
	categoriesTable = "categories"
)

func (s *service) GetCategoriesByUser(ctx context.Context, userId int) (domain.Category, error) {
	var emptyCategory domain.Category

	req, err := createGetCategoriesByUserRequest(userId)
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
		return emptyCategory, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	// fmt.Println("BODY:", string(bodyBytes))
	if err != nil {
		fmt.Printf("Failed to read body: %v\n", err)
		return emptyCategory, err
	}

	var categories []domain.Category
	if err := json.Unmarshal(bodyBytes, &categories); err != nil {
		fmt.Printf("Failed to parse JSON: %v\n", err)
		return emptyCategory, err
	}

	return categories[0], nil
}

func createGetCategoriesByUserRequest(userId int) (*http.Request, error) {
	url := os.Getenv("SUPABASE_URL") + "/rest/v1"
	endpoint := fmt.Sprintf("%s/%s?select=*&profile_id=eq.%d", url, categoriesTable, userId)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return nil, err
	}
	return req, nil
}
