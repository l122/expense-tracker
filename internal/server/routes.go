package server

import (
	"encoding/json"
	"log"
	"net/http"

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
	router := http.NewServeMux()
	handler := &Handler{
		Handler: router,
	}

	templates := web.ParseTemplates()

	router.Handle("GET /", index.NewIndexHandler(index.NewIndexView(templates)))
	router.Handle("GET /navbar", navbar.NewHandler(navbar.NewView(templates)))
	router.Handle("GET /dashboard", dashboard.NewHandler(dashboard.New(templates)))
	router.HandleFunc("GET /health", s.HealthHandler)
	router.Handle("GET /about", about.NewAboutHandler(about.NewAboutView(templates, s.config)))

	// Auth
	router.Handle("GET /login", login.NewLoginHandler(login.NewLoginView(templates), s.db))
	router.Handle("GET /auth/google", login.NewAuthHandler(login.NewLoginView(templates), s.db))
	router.Handle("GET /auth/google/callback", login.NewCallbackHandler(login.NewLoginView(templates), s.db))
	router.Handle("GET /admin", admin.NewAdminHandler(s.db, admin.NewAdminView(templates)))
	router.Handle("PATCH /admin/{id}/enable", admin.NewEnableUserHandler(s.db, admin.NewAdminView(templates)))
	router.Handle("PATCH /admin/{id}/disable", admin.NewDisableUserHandler(s.db, admin.NewAdminView(templates)))

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
