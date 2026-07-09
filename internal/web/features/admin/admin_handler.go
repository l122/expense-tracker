package admin

import (
	"context"
	"net/http"
	"time"

	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/pkgs/mytoken"
	"github.com/l122/expense-tracker/pkgs/redirect"
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

func (t *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	token, err := mytoken.FromRequest(r)
	if err != nil {
		redirect.ToLoginWithError(w, r, "no token in request")
		return
	}
	ctx = mytoken.NewContext(ctx, token)

	users, err := t.repo.GetUsers(ctx)
	if err != nil {
		return
	}

	t.adminView.Index(w, users)
}
