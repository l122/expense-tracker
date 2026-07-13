package navbar

import (
	"context"
	"net/http"
	"time"

	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/pkgs/redirect"
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
	userId, err := userid.FromRequest(r)
	if err != nil {
		// TODO: decide where to redirect
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
	user, err := t.db.GetUserById(ctx, userId)
	if err != nil {
		// TODO: decide where to redirect
		return
	}

	t.view.Index(w, user)
}
