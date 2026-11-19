package models

import "github.com/l122/expense-tracker/internal/domain"

type ChooseCategoryData struct {
	ParentAmount    float64
	ParentNote      string
	Categories      []domain.Category
	ReturnUrl       string
	ParentReturnUrl string
}
