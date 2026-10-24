package categories

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/l122/expense-tracker/internal/domain"
)

func FromHttpResponse(resp *http.Response) ([]domain.Category, error) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read body: %v\n", err)
		return nil, err
	}

	var result []domain.Category
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		fmt.Printf("Failed to parse JSON: %v\n", err)
		return nil, err
	}

	return result, nil
}
