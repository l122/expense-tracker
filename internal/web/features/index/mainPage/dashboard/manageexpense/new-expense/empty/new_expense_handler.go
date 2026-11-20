package empty

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/internal/web/features/index/mainPage/dashboard/manageexpense/models"
	"github.com/l122/expense-tracker/internal/web/features/index/mainPage/dashboard/manageexpense/new-expense/shared"
	"github.com/l122/expense-tracker/pkgs/categories"
	"github.com/l122/expense-tracker/pkgs/httperrors"
	"github.com/l122/expense-tracker/pkgs/token"
	"github.com/l122/expense-tracker/pkgs/userid"
)

type NewExpenseHandler struct {
	http.Handler

	view *shared.NewExpenseView
	repo database.Service
}

func NewNewExpenseHandler(view *shared.NewExpenseView, repo database.Service) *NewExpenseHandler {
	return &NewExpenseHandler{
		view: view,
		repo: repo,
	}
}

func (t *NewExpenseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	categoryId, err := categories.FromUrlQuery(r)
	if err != nil || categoryId == 0 {
		t.view.Index(w, models.ExpensePageData{
			ReturnUrl: "/dashboard",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	accessToken := token.FromRequestOrPanic(r)
	ctx = token.NewContext(ctx, accessToken)

	userId, err := userid.FromRequest(r)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	returnURL := fmt.Sprintf("/category/%d/expenses", categoryId)
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

	data := models.ExpensePageData{
		Category:  category,
		ReturnUrl: returnURL,
	}

	t.view.Index(w, data)
}
