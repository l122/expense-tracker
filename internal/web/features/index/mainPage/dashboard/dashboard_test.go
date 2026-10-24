package dashboard_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/l122/expense-tracker/internal/domain"
	"github.com/l122/expense-tracker/internal/web"
	"github.com/l122/expense-tracker/internal/web/features/index/mainPage/dashboard"
)

func TestDashboardView_WithCategories(t *testing.T) {
	templates := web.ParseTemplates()
	view := dashboard.New(templates)

	categories := []domain.Category{
		{
			Id:        1,
			ProfileId: 10,
			Name:      "Food & Dining",
			Color:     "#FF5733",
			Icon:      "🍔",
		},
		{
			Id:        2,
			ProfileId: 10,
			Name:      "Transportation",
			Color:     "#33A1FF",
			Icon:      "🚗",
		},
	}

	w := httptest.NewRecorder()
	view.Index(w, categories)

	if w.Code != 200 {
		t.Fatalf("expected status code 200, got %d", w.Code)
	}

	body := w.Body.String()
	for _, expected := range []string{"Food &amp; Dining", "🍔", "#FF5733", "Transportation", "🚗", "#33A1FF"} {
		if !strings.Contains(body, expected) {
			t.Errorf("expected body to contain %q, but it did not", expected)
		}
	}
}

func TestDashboardView_EmptyCategories(t *testing.T) {
	templates := web.ParseTemplates()
	view := dashboard.New(templates)

	w := httptest.NewRecorder()
	view.Index(w, []domain.Category{})

	if w.Code != 200 {
		t.Fatalf("expected status code 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "No categories found.") {
		t.Errorf("expected body to contain 'No categories found.', but got %s", body)
	}
}
