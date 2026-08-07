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
	// Check if the request comes from the same user
	userIdFromCookie := userid.FromRequestOrPanic(r)
	userIdFromUrl, err := userid.FromUrlPath(r)
	if err != nil {
		// TODO: log
		fmt.Printf("Error: %v\n", err)
		return
	}
	if userIdFromCookie != userIdFromUrl {
		// TODO: log and redirect to an error page
		// return 403
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	accessToken := token.FromRequestOrPanic(r)
	ctx = token.NewContext(ctx, accessToken)
	err = h.repo.DeleteUsers(ctx, userIdFromUrl)
	if userIdFromCookie != userIdFromUrl {
		// TODO: log and redirect to an error page
		// return 400
		return
	}

	cookies.ClearAll(w, r)

	http.Redirect(
		w,
		r,
		"/auth/login",
		http.StatusTemporaryRedirect,
	)
}
