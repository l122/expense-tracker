package database

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/l122/expense-tracker/internal/domain"
	"github.com/l122/expense-tracker/pkgs/dbhttp"
	"github.com/l122/expense-tracker/pkgs/expenses"
)

func (s *service) GetExpensesByCategoryIdAndUser(ctx context.Context, categoryId int, userId int) ([]domain.Expense, error) {
	var result []domain.Expense
	req, err := createGetExpensesByCategoryIdAndUserRequest(categoryId, userId)
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
		return result, fmt.Errorf("Error response (Status %d): %s\n", resp.StatusCode, string(bodyBytes))
	}

	result, err = expenses.FromHttpResponse(resp)
	if err != nil {
		return result, err
	}

	return result, nil
}

func createGetExpensesByCategoryIdAndUserRequest(categoryId int, userId int) (*http.Request, error) {
	url := os.Getenv("SUPABASE_URL") + "/rest/v1"
	endpoint := fmt.Sprintf("%s/%s?select=*&profile_id=eq.%d&category_id=eq.%d", url, expensesTable, userId, categoryId)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return nil, err
	}
	return req, nil
}
