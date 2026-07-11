package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/l122/expense-tracker/internal/web"
	"github.com/l122/expense-tracker/internal/web/features/index"
	"github.com/l122/expense-tracker/internal/web/features/index/mainPage/about"
	"github.com/l122/expense-tracker/internal/web/features/index/mainPage/admin"
	"github.com/l122/expense-tracker/internal/web/features/index/mainPage/dashboard"
	"github.com/l122/expense-tracker/internal/web/features/login"
	"github.com/l122/expense-tracker/internal/web/features/navbar"
)

type Handler struct {
	http.Handler
}

func (s *Server) RegisterRoutes() *Handler {
	router := mux.NewRouter()
	handler := &Handler{
		Handler: router,
	}

	templates := web.ParseTemplates()

	router.Handle("/", index.NewIndexHandler(index.NewIndexView(templates))).Methods(http.MethodGet)
	router.Handle("/navbar", navbar.NewHandler(navbar.NewView(templates))).Methods(http.MethodGet)
	router.Handle("/dashboard", dashboard.NewHandler(dashboard.New(templates))).Methods(http.MethodGet)
	router.HandleFunc("/health", s.HealthHandler).Methods(http.MethodGet)
	router.Handle("/about", about.NewAboutHandler(about.NewAboutView(templates, s.config))).Methods(http.MethodGet)

	// Auth
	router.Handle("/login", login.NewLoginHandler(login.NewLoginView(templates), s.db)).Methods(http.MethodGet)
	router.Handle("/auth/google", login.NewAuthHandler(login.NewLoginView(templates), s.db)).Methods(http.MethodGet)
	router.Handle("/auth/google/callback", login.NewCallbackHandler(login.NewLoginView(templates), s.db)).Methods(http.MethodGet)
	router.Handle("/admin", admin.NewAdminHandler(s.db, admin.NewAdminView(templates))).Methods(http.MethodGet)
	// router.Handle("/admin/{id}", admin.NewDeleteUserHandler(s.db, admin.NewUsersListView(templates))).Methods(http.MethodDelete)

	return handler
}

func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(s.db.Health())
	if err != nil {
		http.Error(w, "Failed to marshal health check response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
