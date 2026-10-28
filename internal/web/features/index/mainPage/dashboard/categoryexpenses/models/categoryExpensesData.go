package models

import "github.com/l122/expense-tracker/internal/domain"

type CategoryExpensesData struct {
	Category    domain.Category
	TotalAmount float64
	DateGroups  []ExpenseDateGroup
}
