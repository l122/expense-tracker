package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/l122/expense-tracker/cmd/web"
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

	router.Handle("/", web.NewIndexHandler(web.NewIndexView(templates))).Methods(http.MethodGet)

	return handler
}
