package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/l122/expense-tracker/internal/domain"
	"github.com/l122/expense-tracker/pkgs/dbhttp"
)

func (s *service) CreateExpense(ctx context.Context, expense domain.Expense) error {
	req, err := createExpenseRequest(expense)
	if err != nil {
		return err
	}

	resp, err := dbhttp.Send(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("Error response (Status %d): %s\n", resp.StatusCode, string(bodyBytes))
		return fmt.Errorf("Error response (Status %d): %s\n", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func createExpenseRequest(expense domain.Expense) (*http.Request, error) {
	url := os.Getenv("SUPABASE_URL") + "/rest/v1"
	endpoint := fmt.Sprintf("%s/%s", url, expensesTable)

	payload, err := json.Marshal(expense)
	if err != nil {
		fmt.Printf("Failed to marshal expense payload: %v\n", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=minimal")

	return req, nil
}
