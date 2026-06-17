package admin

import (
	"net/http"

	"github.com/l122/expense-tracker/internal/database"
)

type AdminHandler struct {
	http.Handler

	adminView *AdminView
	repo      database.Service
}

func NewAdminHandler(service database.Service, adminView *AdminView) *AdminHandler {
	return &AdminHandler{
		adminView: adminView,
		repo:      service,
	}
}

func (t *AdminHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	users, err := t.repo.GetUsers()
	if err != nil {
		return
	}

	t.adminView.Index(w, users)
}
