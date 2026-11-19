package selectcategory

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/internal/web/features/index/mainPage/dashboard/choosecategory/models"
	"github.com/l122/expense-tracker/pkgs/token"
	"github.com/l122/expense-tracker/pkgs/userid"
)

type ChooseCategoryHandler struct {
	http.Handler

	view *View
	repo database.Service
}

func NewChooseCategoryHandler(view *View, repo database.Service) *ChooseCategoryHandler {
	return &ChooseCategoryHandler{
		view: view,
		repo: repo,
	}
}

func (t *ChooseCategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	data := models.ChooseCategoryData{
		Categories: categories,
		ReturnUrl:  "/expense",
	}

	t.view.Index(w, data)
}
