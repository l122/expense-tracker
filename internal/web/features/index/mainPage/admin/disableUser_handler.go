package admin

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/pkgs/appRole"
	"github.com/l122/expense-tracker/pkgs/redirect"
	"github.com/l122/expense-tracker/pkgs/token"
	"github.com/l122/expense-tracker/pkgs/userid"
)

type DisableUserHandler struct {
	http.Handler

	adminView *AdminView
	repo      database.Service
}

func NewDisableUserHandler(service database.Service, adminView *AdminView) *DisableUserHandler {
	return &DisableUserHandler{
		adminView: adminView,
		repo:      service,
	}
}

func (t *DisableUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	accessToken := token.FromRequestOrPanic(r)
	ctx = token.NewContext(ctx, accessToken)

	userId, err := userid.FromUrlPath(r)
	if err != nil {
		// TODO: log
		fmt.Printf("Error: %v\n", err)
		return
	}
	user, err := t.repo.DisableUser(ctx, userId)
	if err != nil {
		// TODO: log
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println(user)

	users, err := t.repo.GetUsers(ctx)
	if err != nil {
		// TODO: log
		fmt.Printf("Error: %v\n", err)
		return
	}

	t.adminView.Index(w, r, users)
}
