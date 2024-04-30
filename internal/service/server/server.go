package server

import (
	"github.com/kyrare/ya-metrics/internal/application/handlers"
	"github.com/kyrare/ya-metrics/internal/infrastructure/metrics"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Server struct {
	Config  Config
	Storage metrics.Storage
	Logger  zap.SugaredLogger
}

func (s *Server) Run() error {
	if s.Config.Restore {
		err := s.Storage.Restore()

		if err != nil {
			return err
		}
	}

	s.storageStore()

	r := handlers.ServerRouter(s.Storage, s.Config.StoreInterval == 0, s.Logger)

	return http.ListenAndServe(s.Config.Address, r)
}

func (s *Server) storageStore() {
	if s.Config.StoreInterval == 0 {
		return
	}

	ticker := time.NewTicker(s.Config.StoreInterval * time.Second)

	go func() {
		for range ticker.C {
			err := s.Storage.Store()
			if err != nil {
				s.Logger.Error(err)
			}
		}
	}()
}

func NewServer(config Config, storage metrics.Storage, logger zap.SugaredLogger) *Server {
	return &Server{
		Config:  config,
		Storage: storage,
		Logger:  logger,
	}
}
