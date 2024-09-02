package server

import (
	"context"
	"database/sql"

	dm "github.com/kyrare/ya-metrics/internal/domain/metrics"
	"github.com/kyrare/ya-metrics/internal/infrastructure/metrics"
	pb "github.com/kyrare/ya-metrics/internal/infrastructure/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MetricsServer struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedMetricsServer

	storage metrics.Storage
	db      *sql.DB
	logger  zap.SugaredLogger
}

func NewMetricsServer(storage metrics.Storage, db *sql.DB, logger zap.SugaredLogger) *MetricsServer {
	return &MetricsServer{
		storage: storage,
		db:      db,
		logger:  logger,
	}
}

func (s MetricsServer) AddMetrics(ctx context.Context, in *pb.AddMetricsRequest) (*pb.AddMetricsResponse, error) {
	var response pb.AddMetricsResponse

	if len(in.Metrics) == 0 {
		return &response, status.Errorf(codes.InvalidArgument, "metrics required")
	}

	save := make([]dm.Metrics, 0, len(in.Metrics))
	for _, m := range in.Metrics {
		save = append(save, dm.ToGRPC(m))
	}

	s.logger.Info("Обновляем данные по gRPC", save)

	if err := s.storage.Updates(save); err != nil {
		s.logger.Error("Не удалось обновить записи", err, in)
		response.Error = "Не удалось обновить записи"
	}

	return &response, nil
}

func GRPCMetricToMetrics(m *pb.Metric) dm.Metrics {
	var t dm.MetricType

	if m.MType == pb.Metric_GAUGE {
		t = dm.TypeGauge
	}
	if m.MType == pb.Metric_COUNTER {
		t = dm.TypeCounter
	}

	return dm.Metrics{
		ID:    m.Id,
		MType: string(t),
		Delta: &m.Delta,
		Value: &m.Value,
	}
}
