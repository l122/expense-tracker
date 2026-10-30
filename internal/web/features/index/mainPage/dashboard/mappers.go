package dashboard

import (
	"github.com/l122/expense-tracker/internal/domain"
	"github.com/l122/expense-tracker/internal/web/features/index/mainPage/dashboard/models"
)

func newCategoriesList(categories []domain.Category, expenses []domain.Expense) models.CategoriesList {
	var total float64
	for _, e := range expenses {
		total += e.Amount
	}

	return models.CategoriesList{
		Categories:  newCategoryItems(categories, expenses),
		TotalAmount: total,
	}
}

func newCategoryItems(categories []domain.Category, expenses []domain.Expense) []models.CategoryItemData {
	categoryTotals := make(map[int]float64)
	for _, e := range expenses {
		categoryTotals[e.CategoryId] += e.Amount
	}

	var result []models.CategoryItemData
	for _, c := range categories {
		result = append(result, models.CategoryItemData{
			Category: c,
			Amount:   categoryTotals[c.Id],
		})
	}

	return result
}
