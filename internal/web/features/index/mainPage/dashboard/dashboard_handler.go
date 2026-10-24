package dashboard

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/pkgs/token"
	"github.com/l122/expense-tracker/pkgs/userid"
)

type DashboardHandler struct {
	http.Handler

	view *View
	repo database.Service
}

func NewHandler(view *View, repo database.Service) *DashboardHandler {
	return &DashboardHandler{
		view: view,
		repo: repo,
	}
}

func (t *DashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	accessToken := token.FromRequestOrPanic(r)
	ctx = token.NewContext(ctx, accessToken)

	userId, err := userid.FromRequest(r)
	if err != nil {
		// TODO: log
		fmt.Printf("Error: %v\n", err)
		return
	}
	categories, err := t.repo.GetCategoriesByUser(ctx, userId)
	if err != nil {
		// TODO: log
		fmt.Printf("Error fetching categories: %v\n", err)
	}

	t.view.Index(w, categories)
}
