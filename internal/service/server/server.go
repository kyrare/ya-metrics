// Пакет server поднимает api для получения метрики и их хранения

package server

import (
	"database/sql"
	"net"
	"net/http"
	"time"

	"github.com/kyrare/ya-metrics/internal/application/server"
	pb "github.com/kyrare/ya-metrics/internal/infrastructure/proto"

	"github.com/kyrare/ya-metrics/internal/application/handlers"
	"github.com/kyrare/ya-metrics/internal/infrastructure/metrics"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	Config  Config
	Storage metrics.Storage
	DB      *sql.DB
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

	if s.Config.UseGRPC {
		// определяем порт для сервера
		listen, err := net.Listen("tcp", s.Config.Address)
		if err != nil {
			return err
		}
		// создаём gRPC-сервер без зарегистрированной службы
		gRPCServer := grpc.NewServer()
		// регистрируем сервис
		// todo по хорошему тут нужны еще серверы по остальному апи который предоставляется по HTTP,
		// 		но из-за нехватки времени сделал только минимум который требуется для агента
		pb.RegisterMetricsServer(gRPCServer, server.NewMetricsServer(s.Storage, s.DB, s.Logger))

		s.Logger.Info("Start gRPC server")
		// получаем запрос gRPC
		if err := gRPCServer.Serve(listen); err != nil {
			return err
		}
	} else {
		s.Logger.Info("Start HTTP server")

		r := handlers.ServerRouter(s.Storage, s.DB, s.Config.StoreInterval == 0, s.Config.CheckKey, s.Config.CryptoKey, s.Config.TrustedSubnet, s.Logger)

		return http.ListenAndServe(s.Config.Address, r)
	}

	return nil
}

func (s *Server) storageStore() {
	if s.Config.StoreInterval == 0 {
		return
	}

	ticker := time.NewTicker(s.Config.StoreInterval * time.Second)

	s.Logger.Info("Store will be saved by interval, interval - ", s.Config.StoreInterval)

	go func() {
		for range ticker.C {
			err := s.Storage.Store()
			if err != nil {
				s.Logger.Error(err)
			}
		}
	}()
}

func NewServer(config Config, storage metrics.Storage, db *sql.DB, logger zap.SugaredLogger) *Server {
	return &Server{
		Config:  config,
		Storage: storage,
		DB:      db,
		Logger:  logger,
	}
}
