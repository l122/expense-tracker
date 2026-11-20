package new_expense

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/internal/domain"
	"github.com/l122/expense-tracker/pkgs/token"
	"github.com/l122/expense-tracker/pkgs/userid"
)

type CreateExpenseHandler struct {
	http.Handler

	repo database.Service
}

func NewCreateExpenseHandler(repo database.Service) *CreateExpenseHandler {
	return &CreateExpenseHandler{
		repo: repo,
	}
}

func (t *CreateExpenseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	accessToken := token.FromRequestOrPanic(r)
	ctx = token.NewContext(ctx, accessToken)

	userId, err := userid.FromRequest(r)
	if err != nil {
		// TODO: do the same for all other cases
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err = r.ParseForm(); err != nil {
		http.Error(w, "Invalid data form", http.StatusBadRequest)
		return
	}

	categoryIdStr := strings.TrimSpace(r.FormValue("category_id"))
	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil || categoryId <= 0 {
		http.Error(w, "A valid category is required", http.StatusBadRequest)
		return
	}

	amountStr := strings.TrimSpace(r.FormValue("amount"))
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		http.Error(w, "Amount must be greater than 0", http.StatusBadRequest)
		return
	}

	expenseDate := time.Now().UTC()
	dateStr := strings.TrimSpace(r.FormValue("expense_date"))
	if dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			http.Error(w, "Invalid date format (expected YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		expenseDate = parsedDate
	}

	notes := strings.TrimSpace(r.FormValue("notes"))

	err = t.repo.CreateExpense(ctx, domain.Expense{
		ExpenseDate: expenseDate,
		ProfileId:   userId,
		CategoryId:  categoryId,
		Amount:      amount,
		Notes:       notes,
	})
	if err != nil {
		fmt.Printf("Error creating expense: %v\n", err)
		http.Error(w, "Failed to save expense", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
