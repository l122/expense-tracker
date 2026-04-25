package server

import (
	"fmt"
	"net/http"
)

type Server struct {
	port int
}

func NewServer() *http.Server {
	NewServer := &Server{
		port: 8080,
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", NewServer.port),
		Handler: NewServer.RegisterRoutes(),
	}

	return server
}
