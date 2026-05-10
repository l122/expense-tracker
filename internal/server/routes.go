package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/l122/expense-tracker/cmd/web"
	"github.com/l122/expense-tracker/cmd/web/features/about"
	"github.com/l122/expense-tracker/cmd/web/features/admin"
	"github.com/l122/expense-tracker/cmd/web/features/index"
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
	router.Handle("/admin", admin.NewAdminHandler(admin.NewAdminView(templates))).Methods(http.MethodGet)
	router.Handle("/about", about.NewAboutHandler(about.NewAboutView(templates, s.config))).Methods(http.MethodGet)

	return handler
}
