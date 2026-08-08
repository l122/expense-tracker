package admin

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/pkgs/token"
	"github.com/l122/expense-tracker/pkgs/userid"
)

type DeleteUserHandler struct {
	http.Handler

	adminView *AdminView
	repo      database.Service
}

func NewDeleteUserHandler(service database.Service, adminView *AdminView) *DeleteUserHandler {
	return &DeleteUserHandler{
		adminView: adminView,
		repo:      service,
	}
}

func (t *DeleteUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId, err := userid.FromUrlPath(r)
	if err != nil {
		// TODO: log
		fmt.Printf("Error: %v\n", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	accessToken := token.FromRequestOrPanic(r)
	ctx = token.NewContext(ctx, accessToken)
	err = t.repo.DeleteUser(ctx, userId)
	if err != nil {
		// TODO: log
		fmt.Printf("Error: %v\n", err)
	}

	users, err := t.repo.GetUsers(ctx)
	if err != nil {
		// TODO: log
		fmt.Printf("Error: %v\n", err)
		return
	}

	t.adminView.Index(w, r, users)
}
