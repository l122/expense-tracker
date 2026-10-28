package models

import "github.com/l122/expense-tracker/internal/domain"

type ExpenseDateGroup struct {
	DateFormatted string
	Expenses      []domain.Expense
}
