package view_expense

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/internal/web/features/index/mainPage/dashboard/manageexpense/models"
	"github.com/l122/expense-tracker/internal/web/features/index/mainPage/dashboard/manageexpense/new_expense"
	"github.com/l122/expense-tracker/pkgs/expenses"
	"github.com/l122/expense-tracker/pkgs/httperrors"
	"github.com/l122/expense-tracker/pkgs/token"
	"github.com/l122/expense-tracker/pkgs/userid"
)

type ViewExpenseHandler struct {
	http.Handler

	view *new_expense.NewExpenseView
	repo database.Service
}

func NewViewExpenseHandler(view *new_expense.NewExpenseView, repo database.Service) *ViewExpenseHandler {
	return &ViewExpenseHandler{
		view: view,
		repo: repo,
	}
}

func (t *ViewExpenseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	expenseId, err := expenses.FromUrlPath(r)
	if err != nil || expenseId <= 0 {
		http.Error(w, "Error expense id", http.StatusBadRequest)
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

	expense, err := t.repo.GetExpenseByIdForUser(ctx, expenseId, userId)
	returnURL := fmt.Sprintf("/category/%d/expenses", expense.CategoryId)
	if err != nil {
		if errors.Is(err, httperrors.ErrNotFound) {
			fmt.Printf("Expense with id %d and user id %d is not found\n", expenseId, userId)
			http.NotFound(w, r)
		} else {
			fmt.Printf("Error fetching expense: %v\n", err)
			http.Error(w, "Failed to fetch expense", http.StatusInternalServerError)
		}
		return
	}

	category, err := t.repo.GetCategoryByIdForUser(ctx, expense.CategoryId, userId)
	if err != nil {
		if errors.Is(err, httperrors.ErrNotFound) {
			fmt.Printf("Category with id %d and user id %d is not found\n", expense.CategoryId, userId)
			http.NotFound(w, r)
		} else {
			fmt.Printf("Error fetching category: %v\n", err)
			http.Error(w, "Failed to fetch category", http.StatusInternalServerError)
		}
		return
	}

	data := models.ExpensePageData{
		Expense:   expense,
		Category:  category,
		ReturnUrl: returnURL,
	}

	t.view.Index(w, data)
}
