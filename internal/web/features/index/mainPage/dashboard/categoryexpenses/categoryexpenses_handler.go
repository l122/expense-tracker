package categoryexpenses

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/pkgs/categories"
	"github.com/l122/expense-tracker/pkgs/httperrors"
	"github.com/l122/expense-tracker/pkgs/token"
	"github.com/l122/expense-tracker/pkgs/userid"
)

type ExpenseHandler struct {
	http.Handler

	view *View
	repo database.Service
}

func NewHandler(view *View, repo database.Service) *ExpenseHandler {
	return &ExpenseHandler{
		view: view,
		repo: repo,
	}
}

func (t *ExpenseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	accessToken := token.FromRequestOrPanic(r)
	ctx = token.NewContext(ctx, accessToken)

	userId, err := userid.FromRequest(r)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	categoryId, err := categories.FromUrlPath(r)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	category, err := t.repo.GetCategoryByIdForUser(ctx, categoryId, userId)
	if err != nil {
		if errors.Is(err, httperrors.ErrNotFound) {
			fmt.Printf("Category with id %d and user id %d is not found\n", categoryId, userId)
			http.NotFound(w, r)
		} else {
			fmt.Printf("Error fetching category: %v\n", err)
			http.Error(w, "Failed to fetch category", http.StatusInternalServerError)
		}
		return
	}

	expenses, err := t.repo.GetExpensesByCategoryIdAndUser(ctx, categoryId, userId)
	t.view.Index(w, newCategoryExpensesData(category, expenses))
}
