package dashboard

import (
	"net/http"
)

type DashboardHandler struct {
	http.Handler

	view *View
}

func NewHandler(view *View) *DashboardHandler {
	return &DashboardHandler{
		view: view,
	}
}

func (t *DashboardHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	t.view.Index(w)
}
