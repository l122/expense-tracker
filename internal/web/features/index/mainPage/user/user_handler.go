package user

import (
	"net/http"

	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/internal/domain"
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

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	h.view.Index(w, domain.User{
		Id:        1,
		FullName:  "Mino",
		Email:     "minolyndo@gmail.com",
		AvatarUrl: "https://lh3.googleusercontent.com/a/ACg8ocKCygLK4XMNa3M8arfNNOnLbaMQ4PXa94zlKTLKnQOGqqexNR0=s96-c",
	})
}
