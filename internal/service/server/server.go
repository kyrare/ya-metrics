package server

import (
	"github.com/kyrare/ya-metrics/internal/application/handlers"
	"net/http"
)

type Server struct {
	Config Config
}

func (s *Server) Run() error {
	r := handlers.ServerRouter()

	return http.ListenAndServe(s.Config.address, r)
}

func NewServer(config Config) *Server {
	return &Server{
		Config: config,
	}
}
