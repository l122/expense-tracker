package server

import (
	"fmt"
	"net/http"

	"github.com/l122/expense-tracker/internal/config"
	"github.com/l122/expense-tracker/internal/database"
)

type Server struct {
	port   int
	config *config.Config

	db database.Service
}

func NewServer(cfg *config.Config) *http.Server {
	port := config.GetConfigValue("PORT")
	if port == 0 {
		port = 8080
	}

	NewServer := &Server{
		port:   port,
		config: cfg,
		db:     database.New(),
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", NewServer.port),
		Handler: NewServer.RegisterRoutes(),
	}

	return server
}
