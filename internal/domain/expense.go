package domain

import "time"

type Expense struct {
	Id          int       `json:"id"`
	ExpenseDate time.Time `json:"expense_date"`
	ProfileId   int       `json:"profile_id"`
	CategoryId  int       `json:"category_id"`
	Title       string    `json:"title"`
	Amount      float64   `json:"amount"`
	Notes       string    `json:"notes"`
}
