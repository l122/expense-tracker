package models

import "github.com/l122/expense-tracker/internal/domain"

type ExpensePageData struct {
	Expense   domain.Expense
	Category  domain.Category
	ReturnUrl string
}
