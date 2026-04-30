package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/l122/expense-tracker/cmd/web"
	features "github.com/l122/expense-tracker/cmd/web/features/index"
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

	router.Handle("/", features.NewIndexHandler(features.NewIndexView(templates, s.config))).Methods(http.MethodGet)

	return handler
}
