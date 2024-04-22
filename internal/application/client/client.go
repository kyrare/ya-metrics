package client

import (
	"fmt"
	"github.com/kyrare/ya-metrics/internal/domain/metrics"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type Client struct {
	serverAddr string
	Logger     zap.SugaredLogger
}

func (c *Client) Send(metricType metrics.MetricType, metric string, value float64) error {
	uri := fmt.Sprintf("http://%s/update/%s/%s/%v", c.serverAddr, metricType, metric, value)

	c.Logger.Infoln(
		"Начало выполнение запроса",
		"uri", uri,
		"method", "POST",
	)

	response, err := http.Post(uri, "text/plain", nil)

	if err != nil {
		c.Logger.Errorf("Произошла ошибка отправки данных")
		return err
	}

	c.Logger.Infoln(
		"Запрос выполнен",
		"uri", uri,
		"method", "POST",
		"status", response.Status,
		"size", bodySize(response.Body),
	)

	return response.Body.Close()
}

func bodySize(body io.ReadCloser) int {
	bytes, err := io.ReadAll(body)

	if err != nil {
		return 0
	}

	return len(bytes)
}

func NewClient(serverAddr string, logger zap.SugaredLogger) *Client {
	return &Client{serverAddr: serverAddr, Logger: logger}
}
