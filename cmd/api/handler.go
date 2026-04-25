package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	http.Handler

	indexView *IndexView
}

func NewHandler(indexView *IndexView) (*Handler, error) {
	router := mux.NewRouter()
	handler := &Handler{
		Handler: router,

		indexView: indexView,
	}

	router.HandleFunc("/", handler.index).Methods(http.MethodGet)

	return handler, nil
}

func (t *Handler) index(w http.ResponseWriter, _ *http.Request) {
	t.indexView.Index(w)
}
