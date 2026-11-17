package models

import "github.com/l122/expense-tracker/internal/domain"

type ListCategoriesPageData struct {
	Categories []domain.Category
	ReturnUrl  string
}
