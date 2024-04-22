package server

import (
	"github.com/kyrare/ya-metrics/internal/application/handlers"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	Config Config
	Logger zap.SugaredLogger
}

func (s *Server) Run() error {
	r := handlers.ServerRouter(s.Logger)

	return http.ListenAndServe(s.Config.address, r)
}

func NewServer(config Config, logger zap.SugaredLogger) *Server {
	return &Server{
		Config: config,
		Logger: logger,
	}
}
