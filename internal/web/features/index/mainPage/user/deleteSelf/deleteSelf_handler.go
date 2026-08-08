package deleteSelf

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/l122/expense-tracker/internal/database"
	"github.com/l122/expense-tracker/pkgs/cookies"
	"github.com/l122/expense-tracker/pkgs/token"
	"github.com/l122/expense-tracker/pkgs/userid"
)

type Handler struct {
	repo database.Service
}

func New(repo database.Service) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId, err := userid.FromUrlPath(r)
	if err != nil {
		// TODO: log
		fmt.Printf("Error: %v\n", err)
		http.Error(w, "Invalid or missing user id in URL path", http.StatusBadRequest)

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	accessToken := token.FromRequestOrPanic(r)
	ctx = token.NewContext(ctx, accessToken)
	err = h.repo.DeleteUser(ctx, userId)
	if err != nil {
		// TODO: log
		fmt.Printf("Error: %v\n", err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	cookies.ClearAll(w, r)

	// Note:always use this expression for AJAX requests when you need to redirect
	// Otherwise the web page will be embedded
	w.Header().Set("HX-Redirect", "/auth/login")
	w.WriteHeader(http.StatusOK)
}
