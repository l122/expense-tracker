package listcategories

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/internal/web/features/index/mainPage/dashboard/listcategories/models"
	"github.com/l122/expense-tracker/pkgs/token"
	"github.com/l122/expense-tracker/pkgs/userid"
)

type ListCategoriesHandler struct {
	http.Handler

	view *View
	repo database.Service
}

func NewListCategoriesHandler(view *View, repo database.Service) *ListCategoriesHandler {
	return &ListCategoriesHandler{
		view: view,
		repo: repo,
	}
}

func (t *ListCategoriesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	accessToken := token.FromRequestOrPanic(r)
	ctx = token.NewContext(ctx, accessToken)

	userId, err := userid.FromRequest(r)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	categories, err := t.repo.GetCategoriesByUser(ctx, userId)
	if err != nil {
		fmt.Printf("Error fetching categories: %v\n", err)
	}

	data := models.ListCategoriesPageData{
		Categories: categories,
		ReturnUrl:  "/expenses",
	}

	t.view.Index(w, data)
}
