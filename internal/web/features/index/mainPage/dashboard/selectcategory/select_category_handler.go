package selectcategory

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/internal/domain"
	"github.com/l122/expense-tracker/internal/web/features/index/mainPage/dashboard/selectcategory/models"
	"github.com/l122/expense-tracker/pkgs/token"
	"github.com/l122/expense-tracker/pkgs/userid"
)

type ChooseCategoryHandler struct {
	http.Handler

	view *View
	repo database.Service
}

func NewChooseCategoryHandler(view *View, repo database.Service) *ChooseCategoryHandler {
	return &ChooseCategoryHandler{
		view: view,
		repo: repo,
	}
}

func (t *ChooseCategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	accessToken := token.FromRequestOrPanic(r)
	ctx = token.NewContext(ctx, accessToken)

	userId, err := userid.FromRequest(r)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
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

	categories, err := t.repo.GetCategoriesByUser(ctx, userId)
	if err != nil {
		fmt.Printf("Error fetching categories: %v\n", err)
	}

	data := models.ChooseCategoryData{
		Expense: domain.Expense{
			ExpenseDate: expenseDate,
			Amount:      amount,
			Notes:       strings.TrimSpace(r.FormValue("notes")),
		},
		Categories:      categories,
		ReturnUrl:       "/expense",
		ParentReturnUrl: strings.TrimSpace(r.FormValue("return_url")),
	}

	t.view.Index(w, data)
}
