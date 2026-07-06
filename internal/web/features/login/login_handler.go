package login

import (
	"net/http"

	"github.com/l122/expense-tracker/internal/database"
)

type LoginHandler struct {
	view *LoginView
	db   database.Service
}

func NewLoginHandler(view *LoginView, db database.Service) *LoginHandler {
	return &LoginHandler{
		view: view,
		db:   db,
	}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	errParam := r.URL.Query().Get("error")
	data := map[string]interface{}{
		"Error": errParam,
	}
	h.view.Index(w, data)
}
