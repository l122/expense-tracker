package admin

import (
	"net/http"
)

type AdminHandler struct {
	http.Handler

	adminView *AdminView
}

func NewAdminHandler(adminView *AdminView) *AdminHandler {
	return &AdminHandler{
		adminView: adminView,
	}
}

func (t *AdminHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	t.adminView.Index(w)
}
