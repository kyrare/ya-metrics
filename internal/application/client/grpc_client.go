// Package client для отправки данных из агента на сервер
package client

import (
	"context"
	"log"

	"github.com/kyrare/ya-metrics/internal/domain/metrics"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/kyrare/ya-metrics/internal/infrastructure/proto"
)

type GRPCClient struct {
	serverAddr string
	conn       *grpc.ClientConn
	logger     zap.SugaredLogger
}

// NewGRPCClient конструктор для gRPC клиента
func NewGRPCClient(serverAddr string, logger zap.SugaredLogger) (*GRPCClient, error) {
	conn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	c := &GRPCClient{
		serverAddr: serverAddr,
		conn:       conn,
		logger:     logger,
	}

	return c, nil
}

// Send отправляет данные
func (c *GRPCClient) Send(data []metrics.Metrics) {
	send := make([]*pb.Metric, 0, len(data))
	for _, m := range data {
		send = append(send, metrics.FromGRPC(m))
	}

	cl := pb.NewMetricsClient(c.conn)

	resp, err := cl.AddMetrics(context.Background(), &pb.AddMetricsRequest{Metrics: send})

	c.logger.Info("Sended data by gRPC, resp: ", resp)

	if err != nil {
		c.logger.Error(err)
	}
	if resp.Error != "" {
		c.logger.Error(resp.Error)
	}
}

func (c *GRPCClient) Close() error {
	c.logger.Info("CLOSE IN GRPC DELETE ME")
	c.conn.Close()
	return nil
}

func metricsToGRPCMetric(m metrics.Metrics) *pb.Metric {
	var t pb.Metric_MType
	if m.MType == string(metrics.TypeGauge) {
		t = pb.Metric_GAUGE
	}
	if m.MType == string(metrics.TypeCounter) {
		t = pb.Metric_COUNTER
	}

	var delta int64
	if m.Delta != nil {
		delta = *m.Delta
	}

	var value float64
	if m.Value != nil {
		value = *m.Value
	}

	return &pb.Metric{
		Id:    m.ID,
		MType: t,
		Delta: delta,
		Value: value,
	}
}
