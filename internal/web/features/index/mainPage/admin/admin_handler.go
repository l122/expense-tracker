package admin

import (
	"context"
	"net/http"
	"time"

	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/pkgs/appRole"
	"github.com/l122/expense-tracker/pkgs/redirect"
	"github.com/l122/expense-tracker/pkgs/token"
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
	// Check role
	appRole, err := appRole.FromRequest(r)
	if err != nil {
		redirect.ToLoginWithError(w, r, "no AppRole in request")
		return
	}

	if appRole != "admin" {
		// TODO: log and redirect to an error page
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	accessToken, err := token.FromRequest(r)
	if err != nil {
		redirect.ToLoginWithError(w, r, "no token in request")
		return
	}
	ctx = token.NewContext(ctx, accessToken)

	users, err := t.repo.GetUsers(ctx)
	if err != nil {
		// TODO: log
	}

	t.adminView.Index(w, users)
}
