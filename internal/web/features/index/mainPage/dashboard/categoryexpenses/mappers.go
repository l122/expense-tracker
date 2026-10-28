package categoryexpenses

import (
	"slices"
	"time"

	"github.com/l122/expense-tracker/internal/domain"
	"github.com/l122/expense-tracker/internal/web/features/index/mainPage/dashboard/categoryexpenses/models"
	"github.com/l122/expense-tracker/pkgs/utils"
)

func newCategoryExpensesData(category domain.Category, expenses []domain.Expense) models.CategoryExpensesData {
	var total float64
	for _, e := range expenses {
		total += e.Amount
	}

	return models.CategoryExpensesData{
		Category:    category,
		TotalAmount: total,
		DateGroups:  groupByDate(expenses),
	}
}

func groupByDate(expenses []domain.Expense) []models.ExpenseDateGroup {
	grouped := utils.GroupBy(expenses, func(e domain.Expense) time.Time {
		return e.ExpenseDate.Truncate(24 * time.Hour)
	})

	dates := make([]time.Time, 0, len(grouped))
	for d := range grouped {
		dates = append(dates, d)
	}

	slices.SortFunc(dates, func(a, b time.Time) int {
		return b.Compare(a)
	})

	dateGroups := make([]models.ExpenseDateGroup, 0, len(dates))
	for _, d := range dates {
		dateGroups = append(dateGroups, models.ExpenseDateGroup{
			DateFormatted: d.Format("January 02, 2006 (Mon)"),
			Expenses:      grouped[d],
		})
	}

	return dateGroups
}
