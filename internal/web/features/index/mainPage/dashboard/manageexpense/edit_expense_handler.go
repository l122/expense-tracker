package manageexpense

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/internal/domain"
	"github.com/l122/expense-tracker/internal/web/features/index/mainPage/dashboard/manageexpense/models"
	"github.com/l122/expense-tracker/pkgs/categories"
	"github.com/l122/expense-tracker/pkgs/httperrors"
	"github.com/l122/expense-tracker/pkgs/token"
	"github.com/l122/expense-tracker/pkgs/userid"
)

type EditExpenseHandler struct {
	http.Handler

	view *View
	repo database.Service
}

func NewEditExpenseHandler(view *View, repo database.Service) *EditExpenseHandler {
	return &EditExpenseHandler{
		view: view,
		repo: repo,
	}
}

func (t *EditExpenseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	if err != nil || categoryId == 0 {
		t.view.Index(w, models.ExpensePageData{
			ReturnUrl: "/dashboard",
		})
		return
	}

	if err = r.ParseForm(); err != nil {
		http.Error(w, "Invalid data form", http.StatusBadRequest)
		return
	}

	amountStr := strings.TrimSpace(r.FormValue("amount"))
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		amount = 0
	}

	expenseDate := time.Now().UTC()
	dateStr := strings.TrimSpace(r.FormValue("expense_date"))
	if dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			expenseDate = parsedDate
		}
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

	data := models.ExpensePageData{
		Expense: domain.Expense{
			ExpenseDate: expenseDate,
			Amount:      amount,
			Notes:       strings.TrimSpace(r.FormValue("notes")),
		},
		Category:  category,
		ReturnUrl: strings.TrimSpace(r.FormValue("return_url")),
	}

	t.view.Index(w, data)
}
