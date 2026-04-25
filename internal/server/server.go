package server

import (
	"fmt"
	"net/http"

	"github.com/l122/expense-tracker/internal/config"
)

type Server struct {
	port   int
	config *config.Config
}

func NewServer(cfg *config.Config) *http.Server {
	port := config.GetConfigValue("PORT")
	if port == 0 {
		port = 8080
	}

	NewServer := &Server{
		port:   port,
		config: cfg,
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", NewServer.port),
		Handler: NewServer.RegisterRoutes(),
	}

	return server
}
