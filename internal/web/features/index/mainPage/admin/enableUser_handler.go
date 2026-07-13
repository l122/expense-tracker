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
)

type EnableUserHandler struct {
	http.Handler

	adminView *AdminView
	repo      database.Service
}

func NewEnableUserHandler(service database.Service, adminView *AdminView) *EnableUserHandler {
	return &EnableUserHandler{
		adminView: adminView,
		repo:      service,
	}
}

func (t *EnableUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check role
	appRole, err := appRole.FromRequest(r)
	if err != nil {
		redirect.ToLoginWithError(w, r, "no app_role in request")
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

	userId := r.PathValue("id")
	user, err := t.repo.EnableUser(ctx, userId)
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

	t.adminView.Index(w, users)
}
