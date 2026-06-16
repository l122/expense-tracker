package admin

import (
	"net/http"

	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/internal/domain"
)

var users = []domain.User{
	{
		Index:    1,
		Name:     "Alice Johnson",
		Email:    "alice.johnson@example.com",
		Username: "usr_982341",
		Role:     "admin",
	},
	{
		Index:    2,
		Name:     "Bob Smith",
		Email:    "bob.smith@provider.net",
		Username: "usr_441209",
		Role:     "user",
	},
	{
		Index:    3,
		Name:     "Charlie Davis",
		Email:    "charlie.d@service.org",
		Username: "usr_112233",
		Role:     "user",
	},
}

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
