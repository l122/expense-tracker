package index

import (
	"net/http"
)

type IndexHandler struct {
	http.Handler

	indexView *IndexView
}

func NewIndexHandler(indexView *IndexView) *IndexHandler {
	return &IndexHandler{
		indexView: indexView,
	}
}

func (t *IndexHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	t.indexView.Index(w)
}
