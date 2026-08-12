package navbar

import (
	"context"
	"net/http"
	"time"

	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/pkgs/token"
	"github.com/l122/expense-tracker/pkgs/userid"
)

type NavbarHandler struct {
	db database.Service

	view *View
}

func NewHandler(view *View, db database.Service) *NavbarHandler {
	return &NavbarHandler{
		view: view,
		db:   db,
	}
}

func (t *NavbarHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	accessToken := token.FromRequestOrPanic(r)
	ctx = token.NewContext(ctx, accessToken)
	userId := userid.FromRequestOrPanic(r)
	user, err := t.db.GetUserById(ctx, userId)
	if err != nil {
		// TODO: decide where to redirect
		return
	}

	t.view.Index(w, user)
}
