package empty

import (
	"net/http"

	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/internal/web/features/index/mainPage/dashboard/manageexpense/models"
)

type NewExpenseHandler struct {
	http.Handler

	view *NewExpenseView
	repo database.Service
}

func NewNewExpenseHandler(view *NewExpenseView, repo database.Service) *NewExpenseHandler {
	return &NewExpenseHandler{
		view: view,
		repo: repo,
	}
}

func (t *NewExpenseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := models.ExpensePageData{
		ReturnUrl: "/dashboard",
	}

	t.view.Index(w, data)
}
