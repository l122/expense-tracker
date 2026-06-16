package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/l122/expense-tracker/internal/web"
	"github.com/l122/expense-tracker/internal/web/features/about"
	"github.com/l122/expense-tracker/internal/web/features/admin"
	"github.com/l122/expense-tracker/internal/web/features/index"
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
	router.HandleFunc("/health", s.HealthHandler).Methods(http.MethodGet)
	router.Handle("/admin", admin.NewAdminHandler(admin.NewAdminView(templates))).Methods(http.MethodGet)
	router.Handle("/about", about.NewAboutHandler(about.NewAboutView(templates, s.config))).Methods(http.MethodGet)

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
