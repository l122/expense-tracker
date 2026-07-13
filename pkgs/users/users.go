package users

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/l122/expense-tracker/internal/domain"
)

func FromHttpResponse(resp *http.Response, emptyUser domain.User) ([]domain.User, error) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read body: %v\n", err)
		return nil, err
	}

	var users []domain.User
	if err := json.Unmarshal(bodyBytes, &users); err != nil {
		fmt.Printf("Failed to parse JSON: %v\n", err)
		return nil, err
	}
	return users, nil
}
