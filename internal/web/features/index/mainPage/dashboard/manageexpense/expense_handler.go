package expense

import (
	"net/http"

	"github.com/l122/expense-tracker/internal/database"
)

type ManageExpenseHandler struct {
	http.Handler

	view *View
	repo database.Service
}

func New(view *View, repo database.Service) *ManageExpenseHandler {
	return &ManageExpenseHandler{
		view: view,
		repo: repo,
	}
}

func (t *ManageExpenseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.view.Index(w, nil)
}
