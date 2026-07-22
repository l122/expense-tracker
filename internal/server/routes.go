package server

import (
	"net/http"

	"github.com/l122/expense-tracker/internal/web"
	"github.com/l122/expense-tracker/internal/web/features/index"
	"github.com/l122/expense-tracker/internal/web/features/index/mainPage/about"
	"github.com/l122/expense-tracker/internal/web/features/index/mainPage/admin"
	"github.com/l122/expense-tracker/internal/web/features/index/mainPage/dashboard"
	"github.com/l122/expense-tracker/internal/web/features/index/mainPage/user"
	"github.com/l122/expense-tracker/internal/web/features/login"
	"github.com/l122/expense-tracker/internal/web/features/logout"
	"github.com/l122/expense-tracker/internal/web/features/navbar"
)

type Handler struct {
	http.Handler
}

func (s *Server) RegisterRoutes() http.Handler {
	router := http.NewServeMux()

	router.Handle("/", s.registerUserRoutes())
	router.Handle("/auth/", s.registerAuthRoutes())
	router.Handle("/admin/", s.registerAdminRoutes())

	return http.NewCrossOriginProtection().Handler(router)
}

func (s *Server) registerAuthRoutes() http.Handler {
	router := http.NewServeMux()
	handler := &Handler{
		Handler: router,
	}

	templates := web.ParseTemplates()

	router.Handle("GET /auth/logout", logout.New())
	router.Handle("GET /auth/login", login.NewLoginHandler(login.NewLoginView(templates), s.db))
	router.Handle("GET /auth/google", login.NewAuthHandler(login.NewLoginView(templates), s.db))
	router.Handle("GET /auth/google/callback", login.NewCallbackHandler(login.NewLoginView(templates), s.db))

	return chainMiddlewares(handler, rateLimitMiddleware, recoveryMiddleware)
}

func (s *Server) registerUserRoutes() http.Handler {
	router := http.NewServeMux()
	handler := &Handler{
		Handler: router,
	}

	templates := web.ParseTemplates()

	router.Handle("/", index.NewIndexHandler(index.NewIndexView(templates)))
	router.Handle("GET /navbar", navbar.NewHandler(navbar.NewView(templates), s.db))
	router.Handle("GET /dashboard", dashboard.NewHandler(dashboard.New(templates)))
	router.Handle("GET /about", about.NewAboutHandler(about.NewAboutView(templates, s.config)))

	return chainMiddlewares(handler, authMiddleware, rateLimitMiddleware, userRoleCheckMiddleware, recoveryMiddleware)
}

func (s *Server) registerAdminRoutes() http.Handler {
	router := http.NewServeMux()
	handler := &Handler{
		Handler: router,
	}

	templates := web.ParseTemplates()

	router.Handle("GET /admin/", admin.NewAdminHandler(s.db, admin.NewAdminView(templates)))
	router.Handle("GET /admin/{id}", user.NewHandler(s.db, user.NewView(templates)))
	router.Handle("PATCH /admin/{id}/enable", admin.NewEnableUserHandler(s.db, admin.NewAdminView(templates)))
	router.Handle("PATCH /admin/{id}/disable", admin.NewDisableUserHandler(s.db, admin.NewAdminView(templates)))

	return chainMiddlewares(handler, authMiddleware, rateLimitMiddleware, adminRoleCheckMiddleware, recoveryMiddleware)
}
