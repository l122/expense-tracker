package models

import "github.com/l122/expense-tracker/internal/domain"

type ChooseCategoryData struct {
	Categories      []domain.Category
	Expense         domain.Expense
	ReturnUrl       string
	ParentReturnUrl string
}
