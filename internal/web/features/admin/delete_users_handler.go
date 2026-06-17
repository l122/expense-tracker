package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/l122/expense-tracker/internal/database"
)

type DeleteUserHandler struct {
	http.Handler

	adminView *AdminView
	repo      database.Service
}

func NewDeleteUserHandler(service database.Service, adminView *AdminView) *DeleteUserHandler {
	return &DeleteUserHandler{
		adminView: adminView,
		repo:      service,
	}
}

func (t *DeleteUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	fmt.Printf("Deleting user: %d\n", id)

	// TODO: Implement t.repo.DeleteUser(id) here
	if err = t.repo.DeleteUsers(id); err != nil {
		http.Error(w, "Error Deleting user", http.StatusInternalServerError)
		return
	}

	users, err := t.repo.GetUsers()
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}

	t.adminView.Index(w, users)
}
