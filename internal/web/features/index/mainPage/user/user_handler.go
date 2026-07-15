package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/pkgs/redirect"
	"github.com/l122/expense-tracker/pkgs/token"
	"github.com/l122/expense-tracker/pkgs/userid"
)

type UserHandler struct {
	http.Handler

	view *UserView
	repo database.Service
}

func NewHandler(repo database.Service, view *UserView) *UserHandler {
	return &UserHandler{
		view: view,
		repo: repo,
	}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	accessToken, err := token.FromRequest(r)
	if err != nil {
		redirect.ToLoginWithError(w, r, "no token in request")
		return
	}
	ctx = token.NewContext(ctx, accessToken)

	userId, err := userid.FromUrlPath(r)
	if err != nil {
		// TODO: log
		fmt.Printf("Error: %v\n", err)
		return
	}
	user, err := h.repo.GetUserById(ctx, userId)
	if err != nil {
		// TODO: log
		fmt.Printf("User not found: %v\n", err)
	}

	h.view.Index(w, user)
}
