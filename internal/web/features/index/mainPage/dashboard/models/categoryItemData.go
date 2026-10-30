package models

import "github.com/l122/expense-tracker/internal/domain"

type CategoryItemData struct {
	Category domain.Category
	Amount   float64
}
